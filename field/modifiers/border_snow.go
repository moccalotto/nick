package modifiers

import (
	"github.com/moccalotto/nick/field"
	"math/rand"
)

type BorderSnow struct {
	Coverage  float64
	Thickness int
	Alive     bool
	rng       *rand.Rand
}

func NewBorderSnow(Coverage float64, rng *rand.Rand) *BorderSnow {
	return &BorderSnow{Coverage, 1, true, rng}
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

			if inDrawArea && (b.Coverage == 1.0 || b.rng.Float64() < b.Coverage) {
				f.SetAlive(x, y, b.Alive)
			}
		}
	}
}
