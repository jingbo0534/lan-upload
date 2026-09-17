package app
import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	wailsrt "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App Wails 绑定对象
type App struct {
	ctx    context.Context
	cfg    settings
	srv    *ReceiverServer
	ip     string
	saveOK bool
}

func NewApp() *App {
	return &App{cfg: defaultSettings()}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.cfg.load()
	if err := os.MkdirAll(a.cfg.uploadsRoot(), 0755); err != nil {
		a.saveOK = false
	} else {
		a.saveOK = true
	}
	ips := availableIPs()
	if len(ips) > 0 {
		// 优先使用已保存的地址，否则取第一个
		a.ip = ips[0]
		if a.cfg.IP != "" && containsStr(ips, a.cfg.IP) {
			a.ip = a.cfg.IP
		}
	}
	// 若已有配置且未指定，自动启动局域网服务
	a.srv = newReceiverServer(a)
	if a.saveOK {
		if err := a.StartServer(); err != nil {
			fmt.Println("auto start failed:", err)
		}
	}
}

// ---------- 服务器 ----------

func (a *App) GetStatus() ServerStatus {
	url := ""
	port := a.cfg.Port
	if a.srv != nil && a.srv.Running() {
		port = a.srv.Port()
	}
	ips := availableIPs()
	if a.ip != "" && containsStr(ips, a.ip) {
		// keep
	} else if len(ips) > 0 {
		a.ip = ips[0]
	}
	if a.srv != nil && a.srv.Running() && a.ip != "" {
		url = fmt.Sprintf("http://%s:%d/", a.ip, port)
	}
	qr := ""
	if url != "" {
		qr, _ = qrDataURI(url)
	}
	return ServerStatus{
		Running:   a.srv != nil && a.srv.Running(),
		IPs:       ips,
		IP:        a.ip,
		Port:      port,
		URL:       url,
		QrDataURI: qr,
		SaveDir:   a.cfg.SaveDir,
		UploadSub: "uploads",
	}
}

func (a *App) AvailableIPs() []string {
	return availableIPs()
}

func (a *App) StartServer() error {
	if a.saveOK == false {
		return fmt.Errorf("保存目录不可用，请先设置")
	}
	if a.srv.Running() {
		return nil
	}
	if err := a.srv.Start(a.cfg.Port); err != nil {
		return err
	}
	// 记录实际绑定端口
	a.cfg.Port = a.srv.Port()
	_ = a.cfg.save()
	a.notifyServerChanged()
	return nil
}

func (a *App) StopServer() error {
	if a.srv == nil {
		return nil
	}
	return a.srv.Stop()
}

func (a *App) SetServer(ip string, port int) (ServerStatus, error) {
	a.ip = ip
	a.cfg.IP = ip
	wasRunning := a.srv.Running()
	if wasRunning {
		_ = a.srv.Stop()
	}
	a.cfg.Port = port
	if wasRunning {
		if err := a.srv.Start(port); err != nil {
			// 恢复端口为配置值，回报错误
			return a.GetStatus(), err
		}
	}
	_ = a.cfg.save()
	a.notifyServerChanged()
	return a.GetStatus(), nil
}

func (a *App) onServerError(err error) {
	if err == nil {
		return
	}
	if a.srv != nil {
		a.srv.server = nil
		a.srv.ln = nil
	}
	a.notifyServerChanged()
}

func (a *App) notifyServerChanged() {
	if a.ctx == nil {
		return
	}
	st := a.GetStatus()
	wailsrt.EventsEmit(a.ctx, "server:status", st)
}

// notifyUpload 收到新文件时向前端广播（用于右下角提示）
func (a *App) notifyUpload(saved []string) {
	if a.ctx == nil || len(saved) == 0 {
		return
	}
	wailsrt.EventsEmit(a.ctx, "server:upload", saved)
}

// ---------- 设置 ----------

func (a *App) GetConfig() Config {
	return Config{
		Port:       a.cfg.Port,
		SaveDir:    a.cfg.SaveDir,
		AutoFolder: a.cfg.AutoFolder,
		IP:         a.ip,
		WifiSsid:   a.cfg.WifiSSID,
		WifiPass:   a.cfg.WifiPass,
		WifiSec:    a.cfg.WifiSec,
	}
}

// SaveWifi 持久化 WiFi 卡片信息，下次启动自动填充
func (a *App) SaveWifi(ssid, password, sec string) error {
	a.cfg.WifiSSID = ssid
	a.cfg.WifiPass = password
	a.cfg.WifiSec = sec
	return a.cfg.save()
}

func (a *App) SetPort(port int) (ServerStatus, error) {
	return a.SetServer(a.ip, port)
}

func (a *App) SetSaveDir(path string) error {
	if path == "" {
		return fmt.Errorf("路径为空")
	}
	if err := os.MkdirAll(filepath.Join(path, "uploads"), 0755); err != nil {
		return err
	}
	a.cfg.SaveDir = path
	a.saveOK = true
	if err := a.cfg.save(); err != nil {
		return err
	}
	a.notifyServerChanged()
	return nil
}

func (a *App) ChooseSaveDir() (string, error) {
	dir, err := wailsrt.OpenDirectoryDialog(a.ctx, wailsrt.OpenDialogOptions{
		Title: "选择接收文件的保存目录",
	})
	if err != nil || dir == "" {
		return a.cfg.SaveDir, err
	}
	_ = a.SetSaveDir(dir)
	return a.cfg.SaveDir, nil
}

func (a *App) SetAutoFolder(v bool) error {
	a.cfg.AutoFolder = v
	return a.cfg.save()
}

// ---------- 文件管理 ----------

func (a *App) ListFiles(relPath string) ([]FileEntry, error) {
	return a.listEntries(relPath)
}

func (a *App) DeleteFile(path string) error {
	abs, err := a.resolveInUploads(path)
	if err != nil {
		return err
	}
	trash := filepath.Join(a.cfg.SaveDir, ".trash", "manual")
	_ = os.MkdirAll(trash, 0755)
	_, name := filepath.Split(abs)
	if err := os.Rename(abs, filepath.Join(trash, name)); err != nil {
		return err
	}
	a.notifyServerChanged()
	return nil
}

func (a *App) RenameFile(path, newName string) error {
	abs, err := a.resolveInUploads(path)
	if err != nil {
		return err
	}
	newName = sanitizeName(newName)
	if newName == "" {
		return fmt.Errorf("无效名称")
	}
	parent := filepath.Dir(abs)
	if err := os.Rename(abs, filepath.Join(parent, newName)); err != nil {
		return err
	}
	a.notifyServerChanged()
	return nil
}

func (a *App) NewFolder(parent, name string) error {
	base := a.cfg.uploadsRoot()
	if parent != "" {
		abs, err := a.resolveInUploads(parent)
		if err != nil {
			return err
		}
		base = abs
	}
	name = sanitizeName(name)
	if name == "" {
		return fmt.Errorf("无效名称")
	}
	return os.MkdirAll(filepath.Join(base, name), 0755)
}

func (a *App) OpenInExplorer(relPath string) error {
	abs, err := a.resolveInUploads(relPath)
	if err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		cmd := execCommand("explorer.exe", "/select,", abs)
		_ = cmd.Start()
	}
	return nil
}

// ReadFilePreview 读取文件用于预览；图片返回 dataURI，文本返回文本内容
func (a *App) ReadFilePreview(path string) (FilePreview, error) {
	abs, err := a.resolveInUploads(path)
	if err != nil {
		return FilePreview{}, err
	}
	ext := strings.ToLower(filepath.Ext(abs))
	data, err := os.ReadFile(abs)
	if err != nil {
		return FilePreview{}, err
	}
	mime := mimeTypeFor(ext)
	if isImageExt(ext) {
		if thumb, ok := imagePreview(data); ok {
			return FilePreview{
				Previewable: true,
				Kind:        "image",
				Mime:        "image/jpeg",
				Content:     "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(thumb),
			}, nil
		}
		// 无法解码(如 webp/bmp 变体)，回退原始 base64，仅限制体积
		if len(data) > previewImageMaxRaw {
			return FilePreview{}, fmt.Errorf("图片过大，无法预览")
		}
		return FilePreview{
			Previewable: true,
			Kind:        "image",
			Mime:        mime,
			Content:     "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data),
		}, nil
	}
	if isTextExt(ext) {
		if len(data) > 1<<20 {
			return FilePreview{}, fmt.Errorf("文本过大，无法预览")
		}
		return FilePreview{Previewable: true, Kind: "text", Mime: mime, Content: string(data)}, nil
	}
	// 其他类型不预览
	return FilePreview{Previewable: false, Kind: "none", Mime: mime}, nil
}

func (a *App) DownloadAbs(path string) (string, error) {
	return a.resolveInUploads(path)
}

func mimeTypeFor(ext string) string {
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".pdf":
		return "application/pdf"
	case ".txt", ".md", ".log":
		return "text/plain"
	case ".json", ".xml", ".html", ".htm", ".css", ".js":
		return "text/plain"
	case ".mp4", ".webm", ".mov", ".mkv":
		return "video/mp4"
	case ".mp3", ".wav", ".ogg":
		return "audio/mpeg"
	}
	return "application/octet-stream"
}

func isImageExt(ext string) bool {
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".bmp":
		return true
	}
	return false
}

func isTextExt(ext string) bool {
	switch ext {
	case ".txt", ".md", ".log", ".json", ".xml", ".html", ".htm", ".css", ".js", ".csv", ".go", ".py", ".jsx", ".ts":
		return true
	}
	return false
}