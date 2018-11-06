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
	East
	South
	West
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
		pos = Direction(e.rng.Intn(4) + 2)
	}

	distanceFromEdge := int(e.Radius * 3.0 / 4.0)

	switch pos {
	case North:
		f.SetRadius(
			e.rng.Intn(f.Width()),
			distanceFromEdge,
			e.Radius,
			e.Cell,
		)
	case East:
		f.SetRadius(
			f.Width()-distanceFromEdge-1,
			e.rng.Intn(f.Height()),
			e.Radius,
			e.Cell,
		)
	case South:
		f.SetRadius(
			e.rng.Intn(f.Width()),
			f.Height()-distanceFromEdge-1,
			e.Radius,
			e.Cell,
		)
	case West:
		f.SetRadius(
			distanceFromEdge,
			e.rng.Intn(f.Height()),
			e.Radius,
			e.Cell,
		)
	default:
		panic("This should never happen!")
	}
}
