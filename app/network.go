package app
import (
	"net"
	"sort"
	"strings"
)

// isPrivateIPv4 判断是否为局域网私有 IPv4
func isPrivateIPv4(ip net.IP) bool {
	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}
	switch {
	case ip4[0] == 10:
		return true
	case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
		return true
	case ip4[0] == 192 && ip4[1] == 168:
		return true
	case ip4[0] == 169 && ip4[1] == 254: // link-local，跳过
		return false
	}
	return false
}

// availableIPs 枚举本机局域网 IPv4 地址
func availableIPs() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	var out []string
	for _, a := range addrs {
		var ip net.IP
		switch t := a.(type) {
		case *net.IPNet:
			ip = t.IP
		case *net.IPAddr:
			ip = t.IP
		}
		if isPrivateIPv4(ip) {
			out = append(out, ip.String())
		}
	}
	// 稳定排序，优先 192.168
	sort.Slice(out, func(i, j int) bool {
		ii := strings.HasPrefix(out[i], "192.168.")
		jj := strings.HasPrefix(out[j], "192.168.")
		if ii != jj {
			return ii
		}
		return out[i] < out[j]
	})
	return out
}