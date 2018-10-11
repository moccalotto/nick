package printers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/moccaloto/nick/field"
	"image"
	"image/color"
	"image/png"
)

// String returns the game board as a string.
func ItermImage(f *field.Field) string {
	img := image.NewGray(image.Rect(0, 0, f.Width(), f.Height()))

	w := f.Width()
	h := f.Height()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if f.Alive(x, y) {
				img.SetGray(x, y, color.Gray{255})
			}
		}
	}

	buf := new(bytes.Buffer)
	png.Encode(buf, img)

	return fmt.Sprintf(
		"\033]1337;File=inline=1;height=75%%:%s\a",
		base64.StdEncoding.EncodeToString(buf.Bytes()),
	)
}
