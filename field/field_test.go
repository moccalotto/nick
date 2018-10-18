package field

import (
	"testing"
)

func TestCoordsInRange(t *testing.T) {
	f := NewField(10, 10)

	for i := 0; i < f.Width(); i++ {
		for j := 0; j < f.Height(); j++ {
			if !f.CoordsInRange(i, j) {
				t.Errorf("Coords not in range [%d, %d]", i, j)
			}
		}
	}

	edgeTests := [][2]int{
		[2]int{-1, 0},
		[2]int{10, 0},
		[2]int{-1, -1},
		[2]int{10, -1},
		[2]int{-1, 10},
		[2]int{10, 10},
	}

	for _, p := range edgeTests {
		x, y := p[0], p[1]

		if f.CoordsInRange(x, y) {
			t.Errorf("These coords should NOT be allowed: [%d, %d]", x, y)
		}
	}
}
