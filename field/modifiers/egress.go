package modifiers

import (
	"github.com/moccalotto/nick/field"
	"math/rand"
)

type Direction int

const (
	Random    Direction = 0
	North     Direction = 1
	NorthEast Direction = 2
	East      Direction = 3
	SouthEast Direction = 4
	South     Direction = 5
	SouthWest Direction = 6
	West      Direction = 7
	NorthWest Direction = 8
)

type Egress struct {
	Width    int
	Depth    int
	Position Direction
	Alive    bool
	rng      *rand.Rand
}

func NewEgress(position Direction, width int, rng *rand.Rand) *Egress {
	return &Egress{
		Width:    width,
		Depth:    1,
		Position: position,
		Alive:    false, // by default, an egress consists of empty/dead space.
		rng:      rng,
	}
}

func (e *Egress) ApplyToField(f *field.Field) {
	pos := e.Position

	if pos == Random {
		pos = Direction(e.rng.Intn(8) + 1)
	}

	switch pos {
	case North:
		r := Rect{
			StartX:   (f.Width() - e.Width) / 2,
			StartY:   0,
			EndX:     (f.Width() + e.Width) / 2,
			EndY:     e.Depth - 1,
			Coverage: 1.0,
			Alive:    e.Alive,
		}
		r.ApplyToField(f)
	case South:
		r := Rect{
			StartX:   (f.Width() - e.Width) / 2,
			StartY:   f.Height() - e.Depth,
			EndX:     (f.Width() + e.Width) / 2,
			EndY:     f.Height() - 1,
			Coverage: 1.0,
			Alive:    e.Alive,
		}
		r.ApplyToField(f)
	case East:
		r := Rect{
			StartX:   f.Width() - e.Depth,
			StartY:   (f.Height() - e.Width) / 2,
			EndX:     f.Width() - 1,
			EndY:     (f.Height() + e.Width) / 2,
			Coverage: 1.0,
			Alive:    e.Alive,
		}
		r.ApplyToField(f)
	case West:
		r := Rect{
			StartX:   0,
			StartY:   (f.Height() - e.Width) / 2,
			EndX:     e.Depth - 1,
			EndY:     (f.Height() + e.Width) / 2,
			Coverage: 1.0,
			Alive:    e.Alive,
		}
		r.ApplyToField(f)
	case NorthEast:
		l := (e.Width + 1) / 2
		north := Rect{
			StartX:   f.Width() - l,
			StartY:   0,
			EndX:     f.Width() - 1,
			EndY:     e.Depth - 1,
			Coverage: 1.0,
			Alive:    e.Alive,
		}
		north.ApplyToField(f)
		east := Rect{
			StartX:   f.Width() - e.Depth,
			StartY:   0,
			EndX:     f.Width() - 1,
			EndY:     l - 1,
			Coverage: 1.0,
			Alive:    e.Alive,
		}
		east.ApplyToField(f)
	case NorthWest:
		l := (e.Width + 1) / 2
		north := Rect{
			StartX:   0,
			StartY:   0,
			EndX:     l - 1,
			EndY:     e.Depth - 1,
			Coverage: 1.0,
			Alive:    e.Alive,
		}
		north.ApplyToField(f)
		west := Rect{
			StartX:   0,
			StartY:   0,
			EndX:     e.Depth - 1,
			EndY:     l - 1,
			Coverage: 1.0,
			Alive:    e.Alive,
		}
		west.ApplyToField(f)
	case SouthWest:
		l := (e.Width + 1) / 2
		south := Rect{
			StartX:   0,
			StartY:   f.Height() - e.Depth,
			EndX:     l - 1,
			EndY:     f.Height() - 1,
			Coverage: 1.0,
			Alive:    e.Alive,
		}
		south.ApplyToField(f)
		west := Rect{
			StartX:   0,
			StartY:   f.Height() - l,
			EndX:     e.Depth - 1,
			EndY:     f.Height() - 1,
			Coverage: 1.0,
			Alive:    e.Alive,
		}
		west.ApplyToField(f)
	case SouthEast:
		l := (e.Width + 1) / 2
		south := Rect{
			StartX:   f.Width() - l,
			StartY:   f.Height() - e.Depth,
			EndX:     f.Width() - 1,
			EndY:     f.Height() - 1,
			Coverage: 1.0,
			Alive:    e.Alive,
		}
		south.ApplyToField(f)
		east := Rect{
			StartX:   f.Width() - e.Depth,
			StartY:   f.Height() - l,
			EndX:     f.Width() - 1,
			EndY:     f.Height() - 1,
			Coverage: 1.0,
			Alive:    e.Alive,
		}
		east.ApplyToField(f)
	}
}
