package app
import (
	"embed"
	"fmt"
	"io/fs"
	"mime"
	"net"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed uploads/index.html uploads/style.css uploads/app.js
var uploadAssets embed.FS

// staticNames 上传页允许访问的嵌入式文件
var staticNames = map[string]string{
	"/":        "uploads/index.html",
	"/index":   "uploads/index.html",
	"/index.html": "uploads/index.html",
	"/style.css":  "uploads/style.css",
	"/app.js":     "uploads/app.js",
}

// ReceiverServer 局域网上传接收服务
type ReceiverServer struct {
	app    *App
	server *http.Server
	ln     net.Listener
	port   int
}

func newReceiverServer(a *App) *ReceiverServer {
	return &ReceiverServer{app: a}
}

// Start 在指定端口启动监听
func (rs *ReceiverServer) Start(port int) error {
	if rs.server != nil {
		return fmt.Errorf("服务已在运行")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", rs.handleRoot)
	mux.HandleFunc("/upload", rs.handleUpload)
	mux.HandleFunc("/dl/", rs.handleDownload)

	// 绑定端口，冲突则顺延探测空闲端口
	var listener net.Listener
	for probe := port; probe <= port+1000; probe++ {
		ls, e := net.Listen("tcp", fmt.Sprintf(":%d", probe))
		if e == nil {
			listener = ls
			break
		}
	}
	if listener == nil {
		return fmt.Errorf("端口 %d 失败，未找到可用端口", port)
	}
	rs.port = listener.Addr().(*net.TCPAddr).Port
	rs.ln = listener
	rs.server = &http.Server{Handler: mux}

	go func() {
		e := rs.server.Serve(listener)
		rs.app.onServerError(e)
	}()
	return nil
}

// Stop 停止服务
func (rs *ReceiverServer) Stop() error {
	if rs.server == nil {
		return nil
	}
	err := rs.server.Close()
	rs.app.notifyServerChanged()
	rs.server = nil
	rs.ln = nil
	return err
}

// Port 当前实际端口
func (rs *ReceiverServer) Port() int {
	if rs.ln == nil {
		return 0
	}
	return rs.port
}

// Running 是否在运行
func (rs *ReceiverServer) Running() bool {
	return rs.server != nil
}

// handleRoot 首页与静态资源
func (rs *ReceiverServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	embedPath, ok := staticNames[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	data, err := fs.ReadFile(uploadAssets, embedPath)
	if err != nil {
		http.Error(w, "asset not found", 404)
		return
	}
	if strings.HasSuffix(embedPath, ".html") {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	} else if strings.HasSuffix(embedPath, ".css") {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	} else if strings.HasSuffix(embedPath, ".js") {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	}
	w.Write(data)
}

// handleUpload 接收 multipart 上传
func (rs *ReceiverServer) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	// 限制总大小 4GB
	r.Body = http.MaxBytesReader(w, r.Body, 4<<30)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "invalid multipart: "+err.Error(), 400)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "no files", 400)
		return
	}

	saved := make([]string, 0, len(files))
	fails := make([]string, 0, len(files))
	for _, fh := range files {
		src, err := fh.Open()
		if err != nil {
			fails = append(fails, fh.Filename)
			continue
		}
		dest, err := rs.app.saveUploaded(fh.Filename, src)
		src.Close()
		if err != nil {
			fails = append(fails, fh.Filename+" ("+err.Error()+")")
			continue
		}
		saved = append(saved, dest)
	}
	rs.app.notifyServerChanged()
	if len(saved) > 0 {
		rs.app.notifyUpload(saved)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, `{"ok":%v,"saved":%d,"failed":%d,"savedDetail":%s,"failedDetail":%s}`,
		len(fails) == 0,
		len(saved), len(fails),
		quoteJSON(saved), quoteJSON(fails))
}

// handleDownload 下载已接收文件（相对路径经 app 校验，防目录穿越）
func (rs *ReceiverServer) handleDownload(w http.ResponseWriter, r *http.Request) {
	rel := strings.TrimPrefix(r.URL.Path, "/dl/")
	if rel == "" {
		http.NotFound(w, r)
		return
	}
	abs, err := rs.app.resolveInUploads(rel)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename*=UTF-8''"+urlEncode(filepath.Base(abs)))
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(abs)))
	http.ServeFile(w, r, abs)
}

func quoteJSON(v []string) string {
	mb, _ := jsonMarshal(v)
	return string(mb)
}

// resolveInUploads 将相对路径解析为 uploads 下绝对路径（校验穿越）
func (a *App) resolveInUploads(rel string) (string, error) {
	root := a.cfg.uploadsRoot()
	if root == "" {
		return "", fmt.Errorf("no root")
	}
	clean := filepath.Clean("/" + filepath.FromSlash(rel))
	clean = strings.TrimPrefix(clean, string(filepath.Separator))
	if clean == "." || clean == "" {
		return "", fmt.Errorf("bad path")
	}
	abs := filepath.Join(root, clean)
	if !strings.HasPrefix(abs, root+string(filepath.Separator)) && abs != root {
		return "", fmt.Errorf("traversal")
	}
	return abs, nil
}