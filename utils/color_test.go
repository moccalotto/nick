package utils

import (
	"image/color"
	"testing"
)

func TestColorParsers(t *testing.T) {

	tests := map[string]color.NRGBA{
		"rgba(255, 255, 255, 0xff)": color.NRGBA{255, 255, 255, 255},
		"rgb(123, 0xff, 00)":        color.NRGBA{123, 255, 0, 255},
		"#ff00ff":                   color.NRGBA{255, 00, 255, 255},
	}

	for s, res := range tests {
		col, err := ParseColorString(s)

		if res != col {
			t.Errorf("The color string %s did not parse correctly. Got %#v instead (error: %v)", s, col, err)
		}
	}
}
