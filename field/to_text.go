package field

import (
	"strings"
)

// String returns the game board as a string.
func (f *Field) ToText() string {
	var buf strings.Builder
	for y := 0; y < f.h; y++ {
		for x := 0; x < f.w; x++ {
			b := ' '
			if f.Alive(x, y) {
				b = '█'
			}
			buf.WriteRune(b)
		}
		buf.WriteRune('\n')
	}
	return buf.String()
}
