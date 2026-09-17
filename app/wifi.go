package app
import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

// MakeWifiQR 生成"扫码连接 WiFi"的二维码 dataURI
// security: wpa | wep | open
func (a *App) MakeWifiQR(ssid, password, security string) (string, error) {
	ssid = strings.TrimSpace(ssid)
	if ssid == "" {
		return "", fmt.Errorf("请输入 WiFi 名称")
	}
	vs, vp := escWifi(ssid), escWifi(password)
	var content string
	switch security {
	case "wep":
		content = "WIFI:T:WEP;S:" + vs + ";P:" + vp + ";;"
	case "open":
		content = "WIFI:T:nopass;S:" + vs + ";;"
	default: // wpa / wpa2 -> 使用 WPA 以兼容最广泛设备
		content = "WIFI:T:WPA;S:" + vs + ";P:" + vp + ";;"
	}
	return qrDataURI(content)
}

// escWifi 转义 WiFi 二维码里需要转义的特殊字符
func escWifi(s string) string {
	r := strings.NewReplacer("\\", "\\\\", ";", "\\;", ",", "\\,", ":", "\\:", "\"", "\\\"")
	return r.Replace(s)
}

// CurrentWifiSSID 返回当前连接的 WiFi 名称（仅 Windows，失败返回空串）
func (a *App) CurrentWifiSSID() string {
	if runtime.GOOS != "windows" {
		return ""
	}
	out, err := exec.Command("netsh", "wlan", "show", "interfaces").Output()
	if err != nil {
		return ""
	}
	lineRe := regexp.MustCompile(`^\s*SSID\s*:\s*(.+?)\s*$`)
	for _, line := range strings.Split(string(out), "\n") {
		if m := lineRe.FindStringSubmatch(line); m != nil {
			return m[1]
		}
	}
	return ""
}