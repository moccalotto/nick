package effects

import (
	"github.com/moccalotto/nick/field"
	"math/rand"
)

type Rect struct {
	StartX   int
	StartY   int
	EndX     int
	EndY     int
	Coverage float64
	Cell     field.Cell
	rng      *rand.Rand
}

func NewRect(startX, startY, endX, endY int, rng *rand.Rand) *Rect {
	return &Rect{
		StartX:   startX,
		StartY:   startY,
		EndX:     endX,
		EndY:     endY,
		Coverage: 1.0,
		Cell:     field.LivingCell,
		rng:      rng,
	}
}

func (r *Rect) ApplyToField(f *field.Field) {
	if r.Coverage == 0.0 {
		return
	}

	if r.Coverage == 1.0 {
		for y := r.StartY; y <= r.EndY; y++ {
			for x := r.StartX; x <= r.EndX; x++ {
				_ = f.Set(x, y, r.Cell)
			}
		}
		return
	}

	for y := r.StartY; y <= r.EndY; y++ {
		for x := r.StartX; x <= r.EndX; x++ {
			// TODO - this statement could be optimized if coverage is 1.0
			if r.rng.Float64() < r.Coverage {
				_ = f.Set(x, y, r.Cell)
			}
		}
	}
}
