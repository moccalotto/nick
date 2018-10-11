package modifiers

import (
	"github.com/moccaloto/nick/field"
)

type VLine struct {
	StartX      int
	StartY      int
	Length      int
	Probability float64
	Thickness   int
	Alive       bool
}

func NewVLine(startX, startY, length int) *VLine {
	return &VLine{
		StartX:      startX,
		StartY:      startY,
		Length:      length,
		Probability: 1.0,
		Thickness:   1,
		Alive:       true,
	}
}

func (b *VLine) WithThickness(t int) *VLine {
	tmp := *b
	tmp.Thickness = t

	return &tmp
}

func (b *VLine) WithSnow(probability float64) *VLine {
	tmp := *b
	tmp.Probability = probability

	return &tmp
}

// The snow will now add dead cells instead of living cells
func (b *VLine) Inverted(dead bool) *VLine {
	tmp := *b
	tmp.Alive = !dead

	return &tmp
}

func (b *VLine) ToRect() *Rect {
	r := NewRect(
		b.StartX,
		b.StartY,
		b.StartX+b.Thickness-1,
		b.StartY+b.Length-1,
	)

	r.Alive = b.Alive
	r.Probability = b.Probability

	return r
}

func (b *VLine) ApplyToField(f *field.Field) {
	b.ToRect().ApplyToField(f)
}
