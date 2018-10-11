package printers

import (
	"github.com/moccaloto/nick/field"
	"strings"
)

// String returns the game board as a string.
func Text(f *field.Field) string {
	var buf strings.Builder
	for y := 0; y < f.Height(); y++ {
		for x := 0; x < f.Width(); x++ {
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
