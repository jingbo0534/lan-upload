package app
import (
	"encoding/json"
	"os"
	"path/filepath"
)

// settings 本地持久化配置
type settings struct {
	Port       int    `json:"port"`
	SaveDir    string `json:"saveDir"`
	AutoFolder bool   `json:"autoFolder"` // 是否按日期自动分文件夹
	IP         string `json:"ip"`         // 对外广播的本机局域网地址（空则自动取第一个）

	WifiSSID string `json:"wifiSsid"`
	WifiPass string `json:"wifiPass"`
	WifiSec  string `json:"wifiSec"`
}

func defaultSettings() settings {
	configDir, _ := os.UserConfigDir()
	return settings{
		Port:       3210,
		SaveDir:    filepath.Join(configDir, "lan-upload", "received"),
		AutoFolder: true,
	}
}

func (s *settings) configPath() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "lan-upload", "config.json")
}

func (s *settings) load() {
	raw, err := os.ReadFile(s.configPath())
	if err != nil {
		return
	}
	_ = json.Unmarshal(raw, s)
}

func (s *settings) save() error {
	os.MkdirAll(filepath.Dir(s.configPath()), 0755)
	raw, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.configPath(), raw, 0644)
}

func (s *settings) uploadsRoot() string {
	if s.SaveDir == "" {
		return ""
	}
	return filepath.Join(s.SaveDir, "uploads")
}