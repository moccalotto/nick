package utils

import (
	"fmt"
	"image/color"
	"strings"
)

func ParseColorString(s string) (color color.NRGBA, err error) {
	s = strings.ToLower(s)

	if strings.HasPrefix(s, "#") {
		return ParseHexColorString(s)
	}
	if strings.HasPrefix(s, "0x") {
		return ParseHexColorString(s)
	}
	if strings.HasPrefix(s, "rgba(") {
		return ParseRgbaString(s)
	}
	if strings.HasPrefix(s, "rgb(") {
		return ParseRgbString(s)
	}

	err = fmt.Errorf("Invalid color string: '%s'", s)

	return
}

func ParseHexColorString(s string) (col color.NRGBA, err error) {
	s = strings.ToLower(s)
	col.A = 255

	switch len(s) {
	case 8:
		_, err = fmt.Sscanf(s, "0x%02x%02x%02x", &col.R, &col.G, &col.B)
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &col.R, &col.G, &col.B)
	case 5:
		_, err = fmt.Sscanf(s, "0x%1x%1x%1x", &col.R, &col.G, &col.B)
		if err == nil {
			col.R *= 17
			col.G *= 17
			col.B *= 17
		}
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &col.R, &col.G, &col.B)
		if err == nil {
			col.R *= 17
			col.G *= 17
			col.B *= 17
		}
	default:
		err = fmt.Errorf("Invalid hex color: '%s'", s)
	}
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
