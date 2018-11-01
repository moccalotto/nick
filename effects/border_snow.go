package effects

import (
	"github.com/moccalotto/nick/field"
	"math/rand"
)

type BorderSnow struct {
	Coverage  float64
	Thickness int
	Cell      field.Cell
	rng       *rand.Rand
}

func NewBorderSnow(Coverage float64, rng *rand.Rand) *BorderSnow {
	return &BorderSnow{Coverage, 1, field.OnCell, rng}
}

func (this *BorderSnow) ApplyToField(f *field.Field) {
	w := f.Width()
	h := f.Height()

	bw := w - this.Thickness - 1 // x-position of the east line
	bh := h - this.Thickness - 1 // y-position of the south line

	// cannot use async because it would mess up the order in which
	// the rng is used, thus removing reproducability
	f.Map(func(f *field.Field, x, y int, c field.Cell) field.Cell {
		draw := x < this.Thickness || // west line
			y < this.Thickness || // north line
			x > bw || // east line
			y > bh // south line

		if draw && (this.Coverage == 1.0 || this.rng.Float64() < this.Coverage) {
			return this.Cell
		}

		return c
	})
}
