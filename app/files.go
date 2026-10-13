package app
import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// saveUploaded 将上传的流保存到 uploads 目录（带去重改名/日期归档）
func (a *App) saveUploaded(name string, r io.Reader) (string, error) {
	root := a.cfg.uploadsRoot()
	if root == "" {
		return "", os.ErrInvalid
	}
	dir := root
	if a.cfg.AutoFolder {
		dir = filepath.Join(root, time.Now().Format("2006-01-02"))
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	name = sanitizeName(name)
	dest := filepath.Join(dir, name)
	if _, err := os.Stat(dest); err == nil {
		ext := filepath.Ext(name)
		base := strings.TrimSuffix(name, ext)
		for i := 1; ; i++ {
			candidate := filepath.Join(dir, base+" ("+strconv.Itoa(i)+")"+ext)
			if _, err := os.Stat(candidate); os.IsNotExist(err) {
				dest = candidate
				break
			}
		}
	}
	f, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return "", err
	}
	// 返回相对 uploads 根目录的路径，与文件列表/打开目录使用同一套坐标
	return mustRel(root, dest), nil
}

func sanitizeName(name string) string {
	name = filepath.Base(name)
	name = strings.Map(func(r rune) rune {
		switch r {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|', 0:
			return '_'
		}
		return r
	}, name)
	name = strings.TrimSpace(name)
	if name == "" || name == "." {
		return "unnamed"
	}
	return name
}

func mustRel(root, abs string) string {
	rel, _ := filepath.Rel(root, abs)
	return filepath.ToSlash(rel)
}

// listEntries 列出 uploads 下某一目录（relPath 为空表示根）
func (a *App) listEntries(relPath string) ([]FileEntry, error) {
	root := a.cfg.uploadsRoot()
	if root == "" {
		return nil, os.ErrInvalid
	}
	dir := root
	if relPath != "" {
		abs, err := a.resolveInUploads(relPath)
		if err != nil {
			return nil, err
		}
		fi, err := os.Stat(abs)
		if err != nil {
			return nil, err
		}
		if !fi.IsDir() {
			return nil, os.ErrInvalid
		}
		dir = abs
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})
	out := make([]FileEntry, 0, len(entries))
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		full := filepath.Join(dir, e.Name())
		rel, _ := filepath.Rel(root, full)
		fe := FileEntry{
			Name:    e.Name(),
			Path:    filepath.ToSlash(rel),
			IsDir:   e.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime().Format("2006-01-02 15:04"),
			Ext:     strings.ToLower(filepath.Ext(e.Name())),
		}
		if fe.IsDir {
			fe.Size = a.dirSize(full)
		}
		out = append(out, fe)
	}
	return out, nil
}

func (a *App) dirSize(dir string) int64 {
	var total int64
	_ = filepath.Walk(dir, func(_ string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !fi.IsDir() {
			total += fi.Size()
		}
		return nil
	})
	return total
}

// walkUploads 递归收集 uploads 下所有文件
func (a *App) walkUploads() ([]FileEntry, error) {
	var out []FileEntry
	var collect func(rel string) error
	collect = func(rel string) error {
		sub, err := a.listEntries(rel)
		if err != nil {
			return err
		}
		for _, e := range sub {
			if e.IsDir {
				if err := collect(e.Path); err != nil {
					return err
				}
			} else {
				out = append(out, e)
			}
		}
		return nil
	}
	if err := collect(""); err != nil {
		return nil, err
	}
	return out, nil
}