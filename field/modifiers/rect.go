package modifiers

import (
	"github.com/moccaloto/nick/field"
	"math/rand"
)

type Rect struct {
	StartX      int
	StartY      int
	EndX        int
	EndY        int
	Probability float64
	Alive       bool
}

func NewRect(startX, startY, endX, endY int) *Rect {
	return &Rect{
		StartX:      startX,
		StartY:      startY,
		EndX:        endX,
		EndY:        endY,
		Probability: 1.0,
		Alive:       true,
	}
}

func (r *Rect) WithSnow(probability float64) *Rect {
	tmp := *r
	tmp.Probability = probability

	return &tmp
}

// The snow will now add dead cells instead of living cells
func (r *Rect) Inverted(dead bool) *Rect {
	tmp := *r
	tmp.Alive = !dead

	return &tmp
}

func (r *Rect) ApplyToField(f *field.Field) {
	for y := r.StartY; y <= r.EndY; y++ {
		for x := r.StartX; x <= r.EndX; x++ {
			if r.Probability == 1.0 || rand.Float64() < r.Probability {
				f.Set(x, y, r.Alive)
			}
		}
	}
}
