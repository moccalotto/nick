package field

import "testing"

func TestNewField(t *testing.T) {
	f := NewField(3, 3)
	if f.h != 3 {
		t.Errorf("height != %v", 3)
	}
	if f.w != 3 {
		t.Errorf("width != %v", 3)
	}
	if len(f.s) != 3 {
		t.Errorf("y-axis != %v", 3)
	}

	if len(f.s[0]) != 3 {
		t.Errorf("x-axis != %v", 3)
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if f.s[i][j] == true {
				t.Errorf("Cell %d,%d != %v", i, j, false)
			}
		}
	}
}

func TestSet(t *testing.T) {
	f := NewField(3, 3)

	f.Set(0, 0, true)

	if f.s[0][0] == false {
		t.Errorf("cell 0,0 != %v", true)
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 0 && j == 0 {
				// we already tested (0, 0)
				continue
			}
			if f.s[i][j] == true {
				t.Errorf("Cell %d,%d != %v", i, j, false)
			}
		}
	}
}

func TestAlive(t *testing.T) {
	f := NewField(3, 3)

	f.Set(0, 0, true)

	if f.Alive(0, 0) == false {
		t.Errorf("Alive(0,0) != %v", true)
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 0 && j == 0 {
				// we already tested (0, 0)
				continue
			}
			if f.Alive(i, j) == true {
				t.Errorf("Alive(%d,%d) != %v", i, j, false)
			}
		}
	}
}
