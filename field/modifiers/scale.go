package modifiers

import (
	"github.com/moccalotto/nick/field"
	"math"
)

type Scale struct {
	x, y float64
}

func NewScale(f float64) *Scale {
	return &Scale{f, f}
}

func NewScaleXY(x, y float64) *Scale {
	return &Scale{x, y}
}

// Rain living or dead snow onto the given field.
func (s *Scale) ApplyToField(f *field.Field) {
	nw := int(math.Round(float64(f.Width()) * s.x))
	nh := int(math.Round(float64(f.Height()) * s.y))
	tmp := field.NewField(nw, nh)

	for y := 0; y < nh; y++ {
		_y := int(math.Floor(float64(y) / s.y))
		for x := 0; x < nw; x++ {
			_x := int(math.Floor(float64(x) / s.x))
			if f.Alive(_x, _y) {
				tmp.SetAlive(x, y, true)
			}
		}
	}

	f.SetCells(tmp.Cells())
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
