package modifiers

import (
	"github.com/moccalotto/nick/field"
)

// Randomly bring cells to life.
// Each cell has »Coverage« chance to be born.
// NOTE: cells do not die via this method, they are only brought to Coverage.
type Snow struct {
	Coverage float64
	Alive    bool
}

func NewSnow(p float64) *Snow {
	return &Snow{p, true}
}

// Rain living or dead snow onto the given field.
func (s *Snow) ApplyToField(f *field.Field) {
	r := NewRect(0, 0, f.Width()-1, f.Height()-1)
	r.Coverage = s.Coverage
	r.Alive = s.Alive
	r.ApplyToField(f)
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
