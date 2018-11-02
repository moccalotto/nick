package field

// The methods in this enables a field to behave as an image.

import (
	"image"
	"image/color"
	"log"
)

func (f *Field) ColorModel() color.Model {
	return color.AlphaModel
}

func (f *Field) Bounds() image.Rectangle {
	return image.Rect(0, 0, f.w, f.h)
}

func (f *Field) At(x, y int) color.Color {
	idx := int(f.s[x+y*f.w])

	if int(idx) > len(f.Palette) {
		log.Fatalf("Trying to access element no. %d of the palette. But the palette only contains %d elements", idx, len(f.Palette))
	}

	return f.Palette[idx]
}

func (f *Field) ColorIndexAt(x, y int) uint8 {
	return uint8(f.s[x+y*f.w])
}
