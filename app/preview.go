package app
import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

const (
	previewImageMaxDim = 1600 // 预览缩略图最长边
	previewImageMaxRaw = 40 << 20 // 无法解码的图片最大允许体积
)

// imagePreview 将图片压缩为 JPEG 缩略图(最长边 previewImageMaxDim)，
// 返回 JPEG 字节；若无法解码则返回 nil 让调用方回退到原始 base64。
func imagePreview(data []byte) ([]byte, bool) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, false
	}
	scaled := scaleToMax(img, previewImageMaxDim)
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, scaled, &jpeg.Options{Quality: 85}); err != nil {
		return nil, false
	}
	return buf.Bytes(), true
}

// scaleToMax 等比缩放到最长边不超过 maxDim；若已小于则原样返回
func scaleToMax(src image.Image, maxDim int) image.Image {
	b := src.Bounds()
	sw, sh := b.Dx(), b.Dy()
	if sw <= maxDim && sh <= maxDim {
		return src
	}
	scale := float64(maxDim) / float64(max(sw, sh))
	nw := max(1, int(float64(sw)*scale))
	nh := max(1, int(float64(sh)*scale))
	dst := image.NewRGBA(image.Rect(0, 0, nw, nh))
	for y := 0; y < nh; y++ {
		sy := min(b.Min.Y+sh-1, b.Min.Y+int(float64(y)/scale))
		for x := 0; x < nw; x++ {
			sx := min(b.Min.X+sw-1, b.Min.X+int(float64(x)/scale))
			dst.Set(x, y, src.At(sx, sy))
		}
	}
	return dst
}