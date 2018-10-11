package modifiers

import (
	"github.com/moccaloto/nick/field"
)

type Border struct {
	Thickness int
	Alive     bool
}

func NewBorder() *Border {
	return &Border{1, true}
}

func (b *Border) WithThickness(t int) *Border {
	return &Border{t, b.Alive}
}

// The snow will now add dead cells instead of living cells
func (b *Border) Inverted(dead bool) *Border {
	return &Border{
		b.Thickness,
		!dead,
	}
}

func (b *Border) ApplyToField(f *field.Field) {
	w := f.Width()
	h := f.Height()

	bw := w - b.Thickness - 1 // x-position of the east line
	bh := h - b.Thickness - 1 // y-position of the south line

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			draw := x < b.Thickness || // west line
				y < b.Thickness || // north line
				x > bw || // east line
				y > bh // south line

			if draw {
				f.Set(x, y, b.Alive)
			}
		}
	}
}
