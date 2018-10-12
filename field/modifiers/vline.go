package modifiers

import (
	"github.com/moccaloto/nick/field"
)

type VLine struct {
	StartX    int
	StartY    int
	Length    int
	Coverage  float64
	Thickness int
	Alive     bool
}

func NewVLine(startX, startY, length int) *VLine {
	return &VLine{
		StartX:    startX,
		StartY:    startY,
		Length:    length,
		Coverage:  1.0,
		Thickness: 1,
		Alive:     true,
	}
}

func (b *VLine) ToRect() *Rect {
	r := NewRect(
		b.StartX,
		b.StartY,
		b.StartX+b.Thickness-1,
		b.StartY+b.Length-1,
	)

	r.Alive = b.Alive
	r.Coverage = b.Coverage

	return r
}

func (b *VLine) ApplyToField(f *field.Field) {
	b.ToRect().ApplyToField(f)
}
