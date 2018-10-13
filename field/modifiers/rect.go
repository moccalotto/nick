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
	for y := r.StartY; y <= r.EndY; y++ {
		for x := r.StartX; x <= r.EndX; x++ {
			if r.Coverage == 1.0 || rand.Float64() < r.Coverage {
				f.Set(x, y, r.Alive)
			}
		}
	}
}
