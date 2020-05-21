package utils

import (
	"bytes"
	"crypto/md5"
	"image/color"
	"image/png"
)

func rgb(r, g, b uint8) color.NRGBA { return color.NRGBA{r, g, b, 255} }

func md5hash(s string) []byte {
	h := md5.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

func GenIDIcon(rawStr string) []byte {
	var data = md5hash(rawStr) // hex.DecodeString("BM6MqKLq5rBJgHcR6w4p4GXuHSgBuCzxB7LVpRHWP16UTw")
	var config = Sigil{
		Rows: 5,
		Foreground: []color.NRGBA{
			rgb(45, 79, 255),
			rgb(44, 172, 0),
			rgb(254, 180, 44),
			rgb(226, 121, 234),
			rgb(30, 179, 253),
			rgb(232, 77, 65),
			rgb(49, 203, 115),
			rgb(141, 69, 170),
			rgb(252, 125, 31),
		},
		Background: rgb(224, 224, 224),
	}

	img := config.Make(420, false, data)
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
