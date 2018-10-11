package modifiers

import (
	"github.com/moccaloto/nick/field"
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
	Length    int
	Thickness int
	Position  Direction
	Alive     bool
}

func NewEgress(position Direction, length int) *Egress {
	return &Egress{
		Length:    length,
		Thickness: 1,
		Position:  position,
		Alive:     false, // by default, an egress consists of empty/dead space.
	}
}

func (e *Egress) WithThickness(thickness int) *Egress {
	tmp := *e
	tmp.Thickness = thickness

	return &tmp
}

func (e *Egress) Inverted(alive bool) *Egress {
	tmp := *e
	tmp.Alive = alive

	return &tmp
}

func (e *Egress) ApplyToField(f *field.Field) {
	pos := e.Position

	if pos == Random {
		pos = Direction(rand.Intn(8) + 1)
	}

	switch pos {
	case North:
		r := Rect{
			StartX:      (f.Width() - e.Length) / 2,
			StartY:      0,
			EndX:        (f.Width() + e.Length) / 2,
			EndY:        e.Thickness - 1,
			Probability: 1.0,
			Alive:       e.Alive,
		}
		r.ApplyToField(f)
	case South:
		r := Rect{
			StartX:      (f.Width() - e.Length) / 2,
			StartY:      f.Height() - e.Thickness,
			EndX:        (f.Width() + e.Length) / 2,
			EndY:        f.Height() - 1,
			Probability: 1.0,
			Alive:       e.Alive,
		}
		r.ApplyToField(f)
	case East:
		r := Rect{
			StartX:      f.Width() - e.Thickness,
			StartY:      (f.Height() - e.Length) / 2,
			EndX:        f.Width() - 1,
			EndY:        (f.Height() + e.Length) / 2,
			Probability: 1.0,
			Alive:       e.Alive,
		}
		r.ApplyToField(f)
	case West:
		r := Rect{
			StartX:      0,
			StartY:      (f.Height() - e.Length) / 2,
			EndX:        e.Thickness - 1,
			EndY:        (f.Height() + e.Length) / 2,
			Probability: 1.0,
			Alive:       e.Alive,
		}
		r.ApplyToField(f)
	case NorthEast:
		l := (e.Length + 1) / 2
		north := Rect{
			StartX:      f.Width() - l,
			StartY:      0,
			EndX:        f.Width() - 1,
			EndY:        e.Thickness - 1,
			Probability: 1.0,
			Alive:       e.Alive,
		}
		north.ApplyToField(f)
		east := Rect{
			StartX:      f.Width() - e.Thickness,
			StartY:      0,
			EndX:        f.Width() - 1,
			EndY:        l - 1,
			Probability: 1.0,
			Alive:       e.Alive,
		}
		east.ApplyToField(f)
	case NorthWest:
		l := (e.Length + 1) / 2
		north := Rect{
			StartX:      0,
			StartY:      0,
			EndX:        l - 1,
			EndY:        e.Thickness - 1,
			Probability: 1.0,
			Alive:       e.Alive,
		}
		north.ApplyToField(f)
		west := Rect{
			StartX:      0,
			StartY:      0,
			EndX:        e.Thickness - 1,
			EndY:        l - 1,
			Probability: 1.0,
			Alive:       e.Alive,
		}
		west.ApplyToField(f)
	case SouthWest:
		l := (e.Length + 1) / 2
		south := Rect{
			StartX:      0,
			StartY:      f.Height() - e.Thickness,
			EndX:        l - 1,
			EndY:        f.Height() - 1,
			Probability: 1.0,
			Alive:       e.Alive,
		}
		south.ApplyToField(f)
		west := Rect{
			StartX:      0,
			StartY:      f.Height() - l,
			EndX:        e.Thickness - 1,
			EndY:        f.Height() - 1,
			Probability: 1.0,
			Alive:       e.Alive,
		}
		west.ApplyToField(f)
	case SouthEast:
		l := (e.Length + 1) / 2
		south := Rect{
			StartX:      f.Width() - l,
			StartY:      f.Height() - e.Thickness,
			EndX:        f.Width() - 1,
			EndY:        f.Height() - 1,
			Probability: 1.0,
			Alive:       e.Alive,
		}
		south.ApplyToField(f)
		east := Rect{
			StartX:      f.Width() - e.Thickness,
			StartY:      f.Height() - l,
			EndX:        f.Width() - 1,
			EndY:        f.Height() - 1,
			Probability: 1.0,
			Alive:       e.Alive,
		}
		east.ApplyToField(f)
	}
}
