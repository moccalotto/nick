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
	On       bool
	rng      *rand.Rand
}

func NewRect(startX, startY, endX, endY int, rng *rand.Rand) *Rect {
	return &Rect{
		StartX:   startX,
		StartY:   startY,
		EndX:     endX,
		EndY:     endY,
		Coverage: 1.0,
		On:       true,
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
				_ = f.SetOn(x, y, r.On)
			}
		}
		return
	}

	for y := r.StartY; y <= r.EndY; y++ {
		for x := r.StartX; x <= r.EndX; x++ {
			// TODO - this statement could be optimized if coverage is 1.0
			if r.rng.Float64() < r.Coverage {
				_ = f.SetOn(x, y, r.On)
			}
		}
	}
}
