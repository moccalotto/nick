package effects

import (
	"github.com/moccalotto/nick/field"
	"math/rand"
)

type HLine struct {
	StartX    int
	StartY    int
	Length    int
	Coverage  float64
	Thickness int
	On        bool
	rng       *rand.Rand
}

func NewHLine(startX, startY, length int, rng *rand.Rand) *HLine {
	return &HLine{
		StartX:    startX,
		StartY:    startY,
		Length:    length,
		Coverage:  1.0,
		Thickness: 1,
		On:        true,
		rng:       rng,
	}
}

func (b *HLine) Inverted(off bool) *HLine {
	tmp := *b
	tmp.On = !off

	return &tmp
}

func (b *HLine) ToRect() *Rect {
	r := NewRect(
		b.StartX,
		b.StartY,
		b.StartX+b.Length-1,
		b.StartY+b.Thickness-1,
		b.rng,
	)

	r.On = b.On
	r.Coverage = b.Coverage

	return r
}

func (b *HLine) ApplyToField(f *field.Field) {
	b.ToRect().ApplyToField(f)
}
