package app
import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// matchCleanup 按条件判断某文件是否命中
func matchCleanup(a *App, e FileEntry, req CleanupRequest) bool {
	abs, err := a.resolveInUploads(e.Path)
	if err != nil {
		return false
	}
	if len(req.Paths) > 0 {
		return containsPath(req.Paths, e.Path)
	}
	if req.Keyword != "" && !strings.Contains(strings.ToLower(e.Name), strings.ToLower(req.Keyword)) {
		return false
	}
	if len(req.Exts) > 0 && !containsStr(req.Exts, strings.TrimPrefix(e.Ext, ".")) {
		return false
	}
	fi, err := os.Stat(abs)
	if err != nil {
		return false
	}
	if req.Before != "" {
		cut, err := time.ParseInLocation("2006-01-02", req.Before, time.Local)
		if err == nil {
			if !fi.ModTime().Before(cut.Add(24*time.Hour)) {
				return false
			}
		}
	}
	if req.MinSize > 0 && fi.Size() < req.MinSize {
		return false
	}
	if req.MaxSize > 0 && fi.Size() > req.MaxSize {
		return false
	}
	return true
}

// CleanPreview 预览会被清理的文件（供复核）
func (a *App) CleanPreview(req CleanupRequest) (CleanupPreview, error) {
	all, err := a.walkUploads()
	if err != nil {
		return CleanupPreview{}, err
	}
	var out CleanupPreview
	for _, e := range all {
		if matchCleanup(a, e, req) {
			out.Files = append(out.Files, e)
			out.Total++
			out.Size += e.Size
		}
	}
	return out, nil
}

// RunCleanup 将命中文件移入回收站（可恢复），返回结果
func (a *App) RunCleanup(req CleanupRequest) (CleanupResult, error) {
	preview, err := a.CleanPreview(req)
	if err != nil {
		return CleanupResult{}, err
	}
	trash := filepath.Join(a.cfg.SaveDir, ".trash", time.Now().Format("20060102T150405"))
	if err := os.MkdirAll(trash, 0755); err != nil {
		return CleanupResult{}, err
	}
	var res CleanupResult
	for _, e := range preview.Files {
		abs, err := a.resolveInUploads(e.Path)
		if err != nil {
			res.Failed++
			continue
		}
		if err := os.Rename(abs, filepath.Join(trash, e.Name)); err != nil {
			res.Failed++
			continue
		}
		res.Deleted++
	}
	a.notifyServerChanged()
	return res, nil
}

// TrashStats 回收站统计
func (a *App) TrashStats() TrashInfo {
	root := filepath.Join(a.cfg.SaveDir, ".trash")
	files, size := a.walkLocal(root)
	return TrashInfo{Count: int64(len(files)), Size: size}
}

// EmptyTrash 清空回收站（永久删除）
func (a *App) EmptyTrash() error {
	root := filepath.Join(a.cfg.SaveDir, ".trash")
	return os.RemoveAll(root)
}

func (a *App) walkLocal(root string) ([]string, int64) {
	var files []string
	var size int64
	_ = filepath.Walk(root, func(p string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return nil
		}
		files = append(files, p)
		size += fi.Size()
		return nil
	})
	return files, size
}

func containsPath(paths []string, p string) bool {
	for _, x := range paths {
		if x == p {
			return true
		}
	}
	return false
}

func containsStr(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}