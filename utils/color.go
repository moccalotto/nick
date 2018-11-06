package utils

import (
	"fmt"
	"image/color"
	"strings"
)

func ParseColorString(s string) (color color.NRGBA, err error) {
	s = strings.ToLower(s)

	if s[:1] == "#" {
		return ParseHexColorString(s)
	} else if s[:5] == "rgba(" {
		return ParseRgbaString(s)
	} else if s[:4] == "rgb(" {
		return ParseRgbString(s)
	}

	err = fmt.Errorf("Invalid color string: '%s'", s)

	return
}

func ParseHexColorString(s string) (col color.NRGBA, err error) {
	s = strings.ToLower(s)
	col.A = 255

	if len(s) == 7 {
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &col.R, &col.G, &col.B)
		return
	}

	if len(s) == 4 {
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &col.R, &col.G, &col.B)
		if err == nil {
			col.R *= 17
			col.G *= 17
			col.B *= 17
		}
		return
	}

	err = fmt.Errorf("Invalid hex color: '%s'", s)

	return
}

func ParseRgbaString(s string) (col color.NRGBA, err error) {
	_, err = fmt.Sscanf(s, "rgba(%v,%v,%v,%v)", &col.R, &col.G, &col.B, &col.A)

	return
}

func ParseRgbString(s string) (col color.NRGBA, err error) {
	col.A = 255
	_, err = fmt.Sscanf(s, "rgb(%v,%v,%v)", &col.R, &col.G, &col.B)

	return
}
