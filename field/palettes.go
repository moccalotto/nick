package field

import (
	"image/color"
)

var binaryPalette color.Palette
var invertedBinaryPalette color.Palette
var defaultPalette color.Palette

// The default palette
// Off-cells are completely opaque and other cells are of increasing shades of gray
func DefaultPalette() color.Palette {
	if len(defaultPalette) == 0 {
		defaultPalette = make(color.Palette, 256)
		defaultPalette[0] = color.Alpha{0xff}
		opacity := uint8(144)
		for i := 1; i < 255; i++ {
			defaultPalette[i] = color.Alpha{opacity}
			if opacity < 255 {
				opacity += 1
			}
		}
	}
	return defaultPalette
}

// Get the default binary palette
// Off-cells are completely opaque and all other cells are completely transparent.
func BinaryPalette() color.Palette {
	if len(binaryPalette) == 0 {
		binaryPalette = make(color.Palette, 256)
		binaryPalette[0] = color.Alpha{255}
		for i := 1; i < 255; i++ {
			binaryPalette[i] = color.Alpha{0}
		}
	}
	return binaryPalette
}

// Get the default binary palette
// Off-cells are completely opaque and all other cells are completely transparent.
func InvertedBinaryPalette() color.Palette {
	if len(invertedBinaryPalette) == 0 {
		invertedBinaryPalette = make(color.Palette, 256)
		invertedBinaryPalette[0] = color.Alpha{0}
		for i := 1; i < 255; i++ {
			invertedBinaryPalette[i] = color.Alpha{255}
		}
	}
	return invertedBinaryPalette
}
