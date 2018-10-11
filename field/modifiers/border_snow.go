package modifiers

import (
	"github.com/moccaloto/nick/field"
	"math/rand"
)

type BorderSnow struct {
	Probability float64
	Thickness   int
	Alive       bool
}

func NewBorderSnow(probability float64) *BorderSnow {
	return &BorderSnow{probability, 1, true}
}

func (b *BorderSnow) WithThickness(t int) *BorderSnow {
	return &BorderSnow{b.Probability, t, b.Alive}
}

// The snow will now add dead cells instead of living cells
func (b *BorderSnow) Inverted(dead bool) *BorderSnow {
	return &BorderSnow{
		b.Probability,
		b.Thickness,
		!dead,
	}
}

func (b *BorderSnow) ApplyToField(f *field.Field) {
	w := f.Width()
	h := f.Height()

	bw := w - b.Thickness - 1 // x-position of the east line
	bh := h - b.Thickness - 1 // y-position of the south line

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			inDrawArea := x < b.Thickness || // west line
				y < b.Thickness || // north line
				x > bw || // east line
				y > bh // south line

			if inDrawArea && (b.Probability == 1.0 || rand.Float64() < b.Probability) {
				f.Set(x, y, b.Alive)
			}
		}
	}
}
