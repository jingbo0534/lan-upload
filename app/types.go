package app
// ServerStatus 返回给前端的服务器概览信息
type ServerStatus struct {
	Running   bool     `json:"running"`
	IPs       []string `json:"ips"`
	IP        string   `json:"ip"`
	Port      int      `json:"port"`
	URL       string   `json:"url"`
	QrDataURI string   `json:"qrDataUri"`
	SaveDir   string   `json:"saveDir"`
	UploadSub string   `json:"uploadSub"`
}

// FileEntry 文件树中的一个条目（相对路径基于 uploads 根目录）
type FileEntry struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	IsDir   bool   `json:"isDir"`
	Size    int64  `json:"size"`
	ModTime string `json:"modTime"`
	Ext     string `json:"ext"`
}

// FilePreview 文件预览载荷
type FilePreview struct {
	Previewable bool   `json:"previewable"`
	Kind        string `json:"kind"` // image | text | none
	Mime        string `json:"mime"`
	Content     string `json:"content"`
}

// Config 接收端配置
type Config struct {
	Port       int    `json:"port"`
	SaveDir    string `json:"saveDir"`
	AutoFolder bool   `json:"autoFolder"`
	IP         string `json:"ip"`
	WifiSsid   string `json:"wifiSsid"`
	WifiPass   string `json:"wifiPass"`
	WifiSec    string `json:"wifiSec"`
}

// TrashInfo 回收站统计
type TrashInfo struct {
	Count int64 `json:"count"`
	Size  int64 `json:"size"`
}

// CleanupRequest 清理筛选/删除条件
type CleanupRequest struct {
	Paths   []string `json:"paths"`   // 用户勾选的指定文件（相对路径）
	Keyword string   `json:"keyword"` // 名称关键字
	Exts    []string `json:"exts"`    // 指定扩展名
	Before  string   `json:"before"`  // 早于该日期（YYYY-MM-DD）
	MinSize int64    `json:"minSize"` // 大于等于该大小(B)
	MaxSize int64    `json:"maxSize"` // 小于等于该大小(B)
}

// CleanupPreview 清理预览：将命中的文件列给用户复核
type CleanupPreview struct {
	Total int64       `json:"total"`
	Size  int64       `json:"size"`
	Files []FileEntry `json:"files"`
}

// CleanupResult 清理执行结果
type CleanupResult struct {
	Deleted int64 `json:"deleted"`
	Failed  int64 `json:"failed"`
}