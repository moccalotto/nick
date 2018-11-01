package effects

import (
	"github.com/moccalotto/nick/field"
	"math/rand"
)

type Direction int

const (
	Random Direction = iota + 1
	North
	NorthEast
	East
	SouthEast
	South
	SouthWest
	West
	NorthWest
)

type Egress struct {
	Radius   float64
	Position Direction
	On       bool
	rng      *rand.Rand
}

func NewEgress(position Direction, radius float64, rng *rand.Rand) *Egress {
	return &Egress{
		Radius:   radius,
		Position: position,
		On:       false, // by default, an egress consists of empty pace.
		rng:      rng,
	}
}

func (e *Egress) ApplyToField(f *field.Field) {
	pos := e.Position

	if pos == Random {
		pos = Direction(e.rng.Intn(8) + 1)
	}

	cell := field.OffCell
	if e.On {
		cell = field.OnCell
	}

	switch pos {
	case North:
		f.SetRadius(
			f.Width()/2,
			0,
			e.Radius,
			cell,
		)
	case South:
		f.SetRadius(
			f.Width()/2,
			f.Height()-1,
			e.Radius,
			cell,
		)
	case East:
		f.SetRadius(
			f.Width()-1,
			f.Height()/2,
			e.Radius,
			cell,
		)
	case West:
		f.SetRadius(
			0,
			f.Height()/2,
			e.Radius,
			cell,
		)
	case NorthEast:
		f.SetRadius(
			f.Width()-1,
			0,
			e.Radius,
			cell,
		)
	case NorthWest:
		f.SetRadius(
			0,
			0,
			e.Radius,
			cell,
		)
	case SouthWest:
		f.SetRadius(
			0,
			f.Height()-1,
			e.Radius,
			cell,
		)
	case SouthEast:
		f.SetRadius(
			f.Width()-1,
			f.Height()-1,
			e.Radius,
			cell,
		)
	}
}
