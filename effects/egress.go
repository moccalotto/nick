package effects

import (
	"github.com/moccalotto/nick/field"
	"log"
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
	Cell     field.Cell
	rng      *rand.Rand
}

func NewEgress(position Direction, radius float64, rng *rand.Rand) *Egress {

	if radius < 1 {
		log.Fatal("Radius too low")
	}

	return &Egress{
		Radius:   radius,
		Position: position,
		Cell:     field.OffCell,
		rng:      rng,
	}
}

func (e *Egress) ApplyToField(f *field.Field) {
	pos := e.Position

	if pos == Random {
		pos = Direction(e.rng.Intn(8) + 2)
	}

	offset := int(e.Radius / 2.0)

	switch pos {
	case North:
		f.SetRadius(
			f.Width()/2,
			offset,
			e.Radius,
			e.Cell,
		)
	case NorthEast:
		f.SetRadius(
			f.Width()-offset,
			offset,
			e.Radius,
			e.Cell,
		)
	case East:
		f.SetRadius(
			f.Width()-1,
			f.Height()/2,
			e.Radius,
			e.Cell,
		)
	case SouthEast:
		f.SetRadius(
			f.Width()-offset,
			f.Height()-offset,
			e.Radius,
			e.Cell,
		)
	case South:
		f.SetRadius(
			f.Width()/2,
			f.Height()-offset,
			e.Radius,
			e.Cell,
		)
	case SouthWest:
		f.SetRadius(
			offset,
			f.Height()-offset,
			e.Radius,
			e.Cell,
		)
	case West:
		f.SetRadius(
			offset,
			f.Height()/2,
			e.Radius,
			e.Cell,
		)
	case NorthWest:
		f.SetRadius(
			offset,
			offset,
			e.Radius,
			e.Cell,
		)
	default:
		panic("This should never happen!")
	}
}
