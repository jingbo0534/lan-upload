package app
import (
	"bytes"
	"encoding/base64"

	qrcode "github.com/skip2/go-qrcode"
)

// qrDataURI 生成二维码 PNG 并返回 dataURI
func qrDataURI(content string) (string, error) {
	png, err := qrcode.Encode(content, qrcode.Medium, 320)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	b.WriteString("data:image/png;base64,")
	b.WriteString(base64.StdEncoding.EncodeToString(png))
	return b.String(), nil
}