package field

import "math/rand"

// Randomly bring cells to life.
// Each cell has aliveProbability chance to be born.
// NOTE: cells do not die via this method, they are only brought to life.
func (f *Field) Seed(aliveProbability float64) {
	for y := 0; y < f.h; y++ {
		for x := 0; x < f.w; x++ {
			if rand.Float64() < aliveProbability {
				f.Set(x, y, true)
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
