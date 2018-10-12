package modifiers

import (
	"github.com/moccaloto/nick/field"
)

type HLine struct {
	StartX    int
	StartY    int
	Length    int
	Coverage  float64
	Thickness int
	Alive     bool
}

func NewHLine(startX, startY, length int) *HLine {
	return &HLine{
		StartX:    startX,
		StartY:    startY,
		Length:    length,
		Coverage:  1.0,
		Thickness: 1,
		Alive:     true,
	}
}

// The snow will now add dead cells instead of living cells
func (b *HLine) Inverted(dead bool) *HLine {
	tmp := *b
	tmp.Alive = !dead

	return &tmp
}

func (b *HLine) ToRect() *Rect {
	r := NewRect(
		b.StartX,
		b.StartY,
		b.StartX+b.Length-1,
		b.StartY+b.Thickness-1,
	)

	r.Alive = b.Alive
	r.Coverage = b.Coverage

	return r
}

func (b *HLine) ApplyToField(f *field.Field) {
	b.ToRect().ApplyToField(f)
}
