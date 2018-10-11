package modifiers

import (
	"github.com/moccaloto/nick/field"
)

type HLine struct {
	StartX      int
	StartY      int
	Length      int
	Probability float64
	Thickness   int
	Alive       bool
}

func NewHLine(startX, startY, length int) *HLine {
	return &HLine{
		StartX:      startX,
		StartY:      startY,
		Length:      length,
		Probability: 1.0,
		Thickness:   1,
		Alive:       true,
	}
}

func (b *HLine) WithThickness(t int) *HLine {
	tmp := *b
	tmp.Thickness = t

	return &tmp
}

func (b *HLine) WithSnow(probability float64) *HLine {
	tmp := *b
	tmp.Probability = probability

	return &tmp
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
	r.Probability = b.Probability

	return r
}

func (b *HLine) ApplyToField(f *field.Field) {
	b.ToRect().ApplyToField(f)
}
