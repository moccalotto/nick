package modifiers

import (
	"math/rand"
	"github.com/moccaloto/nick/field"
)


// Randomly bring cells to life.
// Each cell has aliveProbability chance to be born.
// NOTE: cells do not die via this method, they are only brought to life.

type Snow struct {
	probability float64
	alive       bool
}

func NewSnow(p float64) *Snow {
	return &Snow{p, true}
}

// The snow will now add dead cells instead of living cells
func (s *Snow) Inverse() *Snow {
	return &Snow{
		s.probability,
		!s.alive,
	}
}

// Rain living or dead snow onto the given field.
func (s *Snow) ApplyToField(f *field.Field) {
	for y := 0; y < f.Height(); y++ {
		for x := 0; x < f.Width(); x++ {
			if rand.Float64() < s.probability {
				f.Set(x, y, s.alive)
			}
		}
	}
}

/* TODO
Add Common shapes such as
Border all the way around
  - Complete
  - Randomly dotted
Lines
  - Complete lines
  - Dotted
  - Dashed
  - Randomly dotted (0 - 100%)
Boxes
  - Border Only
  - Randomly dotted inside (from 0 to 100%)
  - Randomly dotted along border (0 to 100%)
  - Completely Filled (Same as randomly dotted 100%)
Circles
  - Outline Only
  - Randomly dotted inside (from 0 to 100%)
  - Randomly dotted, but with bias towards center
  - Randomly dotted, but with bias towards edge
  - Completely Filled (Same as randomly dotted 100%)
Triangle
Thiele Number Patterns
Exits (in the border)
*/
