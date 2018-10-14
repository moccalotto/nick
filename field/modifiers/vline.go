package modifiers

import (
	"github.com/moccalotto/nick/field"
	"math/rand"
)

type VLine struct {
	StartX    int
	StartY    int
	Length    int
	Coverage  float64
	Thickness int
	Alive     bool
	rng       *rand.Rand
}

func NewVLine(startX, startY, length int, rng *rand.Rand) *VLine {
	return &VLine{
		StartX:    startX,
		StartY:    startY,
		Length:    length,
		Coverage:  1.0,
		Thickness: 1,
		Alive:     true,
		rng:       rng,
	}
}

func (b *VLine) ToRect() *Rect {
	r := NewRect(
		b.StartX,
		b.StartY,
		b.StartX+b.Thickness-1,
		b.StartY+b.Length-1,
		b.rng,
	)

	r.Alive = b.Alive
	r.Coverage = b.Coverage

	return r
}

func (b *VLine) ApplyToField(f *field.Field) {
	b.ToRect().ApplyToField(f)
}
