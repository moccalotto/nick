package effects

import (
	"github.com/moccalotto/nick/field"
)

type Border struct {
	Thickness int
	Cell      field.Cell
}

func NewBorder() *Border {
	return &Border{1, true}
}

func (b *Border) ApplyToField(f *field.Field) {
	w := f.Width()
	h := f.Height()

	bw := w - b.Thickness - 1 // x-position of the east line
	bh := h - b.Thickness - 1 // y-position of the south line

	f.Map(func(f *field.Field, x, y int, c field.Cell) field.Cell {
		draw := x < b.Thickness || // west line
			y < b.Thickness || // north line
			x > bw || // east line
			y > bh // south line

		if !draw {
			return c
		}

		return b.Cell
	})
}
