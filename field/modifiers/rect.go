package modifiers

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
	Alive    bool
}

func NewRect(startX, startY, endX, endY int) *Rect {
	return &Rect{
		StartX:   startX,
		StartY:   startY,
		EndX:     endX,
		EndY:     endY,
		Coverage: 1.0,
		Alive:    true,
	}
}

func (r *Rect) ApplyToField(f *field.Field) {
	if r.Coverage == 0.0 {
		return
	}

	if r.Coverage == 1.0 {
		for y := r.StartY; y <= r.EndY; y++ {
			for x := r.StartX; x <= r.EndX; x++ {
				f.Set(x, y, r.Alive)
			}
		}
		return
	}

	for y := r.StartY; y <= r.EndY; y++ {
		for x := r.StartX; x <= r.EndX; x++ {
			// TODO - this statement could be optimized if coverage is 1.0
			if rand.Float64() < r.Coverage {
				f.Set(x, y, r.Alive)
			}
		}
	}
}
