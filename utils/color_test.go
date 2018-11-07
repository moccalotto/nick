package utils

import (
	"image/color"
	"testing"
)

func TestColorParsers(t *testing.T) {

	tests := map[string]color.NRGBA{
		"rgba(255, 255, 255, 0xff)": color.NRGBA{0xff, 0xff, 0xff, 0xff},
		"rgb(123, 0xff, 00)":        color.NRGBA{0x7b, 0xff, 0x00, 0xff},
		"#ff00ff":                   color.NRGBA{0xff, 0x00, 0xff, 0xff},
		"#fa0":                      color.NRGBA{0xff, 0xaa, 0x00, 0xff},
		"0xff00ff":                  color.NRGBA{0xff, 0x00, 0xff, 0xff},
		"0xfa0":                     color.NRGBA{0xff, 0xaa, 0x00, 0xff},
	}

	for s, res := range tests {
		col, err := ParseColorString(s)

		if res != col {
			t.Errorf("The color string %s did not parse correctly. Got %#v instead (error: %v)", s, col, err)
		}
	}
}
