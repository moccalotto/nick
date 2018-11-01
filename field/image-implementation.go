package field

// The methods in this enables a field to behave as an image.

import (
	"image"
	"image/color"
)

func (f *Field) ColorModel() color.Model {
	return color.AlphaModel
}

func (f *Field) Bounds() image.Rectangle {
	return image.Rect(0, 0, f.w, f.h)
}

func (f *Field) At(x, y int) color.Color {
	idx := f.s[x+y*f.w]

	return f.Palette[idx]
}

func (f *Field) ColorIndexAt(x, y int) uint8 {
	return uint8(f.s[x+y*f.w])
}
