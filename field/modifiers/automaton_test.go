package modifiers

import (
	"github.com/moccalotto/nick/field"
	"testing"
)

func TestNextCellState(t *testing.T) {
	ca := NewAutomaton("B0/S") // rules that inverses
	f := field.NewField(3, 3)

	f.Set(0, 0, true)

	if ca.NextCellState(f, 0, 0) == true {
		t.Errorf("NextCellState(0,0) != %v", false)
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 0 && j == 0 {
				// we already tested (0, 0)
				continue
			}
			if ca.NextCellState(f, i, j) == true {
				t.Errorf("NextCellState(%d,%d) != %v", i, j, true)
			}
		}
	}
}
