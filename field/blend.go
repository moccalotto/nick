package field

import (
	"fmt"
)

const ()

type Blender func(x, y int, source, target Cell) Cell

// Does nothing. Target cell is unchanged.
func BlenderIgnore(x, y int, source, target Cell) Cell {
	return target
}

// Value of target cell becomes value of source cell.
func BlenderCopy(x, y int, source, target Cell) Cell {
	return source
}

// Value of cell in source field is added to the value of the target cell.
func BlenderAdd(x, y int, source, target Cell) Cell {
	return source + target
}

// Value of source cell is subtracted from value of target cell.
func BlenderSub(x, y int, source, target Cell) Cell {
	return source - target
}

// Target cell has its value set to at most 0. Then the value from the source cell is subtracted.
func BlenderMax0ThenSub(x, y int, source, target Cell) Cell {
	if target > 0 {
		target = 0
	}

	return target - source
}

// Source cell is subtracted from target cell. The max value of target cell after subtraction is 0.
func BlenderSubThenMax0(x, y int, source, target Cell) Cell {
	tmp := target - source
	if tmp > 0 {
		tmp = 0
	}
	return tmp
}

// Target cell has its value set to at least 0. Then the source cell's value is added.
func BlenderMin0ThenSub(x, y int, source, target Cell) Cell {
	if target < 0 {
		target = 0
	}

	return target - source
}

// Source cell is added to target cell. Then the target cell's value is adjusted to become at least 0.
func BlenderAddThenMin0(x, y int, source, target Cell) Cell {
	tmp := source + target

	if tmp < 0 {
		return 0
	}
	return tmp
}

// Target cell has its value set to at most 1. Then the value from the source cell is subtracted.
func BlenderMax1ThenSub(x, y int, source, target Cell) Cell {
	if target > 1 {
		target = 1
	}

	return target - source
}

// Source cell is subtracted from target cell. The max value of target cell after subtraction is 1.
func BlenderSubThenMax1(x, y int, source, target Cell) Cell {
	tmp := target - source

	if tmp > 1 {
		return 1
	}

	return tmp
}

// Target cell has its value set to at least 1. THen the source cell's value is added.
func BlenderMin1ThenAdd(x, y int, source, target Cell) Cell {
	if target < 1 {
		target = 1
	}

	return target + source
}

// Source cell is added to target cell. Then the target cell's value is adjusted to become at least 1.
func BlenderAddThenMin1(x, y int, source, target Cell) Cell {
	tmp := source + target

	if tmp < 1 {
		return 1
	}

	return tmp
}

// BLend this field onto a target field.
func (f *Field) BlendOnto(target *Field, alive, dead Blender) error {
	if f.w != target.w || f.h != target.h {
		return fmt.Errorf("Target field has invalid dimensions [%d x %d]. Expected [%d x %d]", target.w, target.h, f.w, f.h)
	}

	for y := 0; y < f.h; y++ {
		for x := 0; x < f.w; x++ {
			idx := y*f.w + x
			targetCell := target.s[idx]
			sourceCell := f.s[idx]
			if sourceCell.Alive() {
				target.s[idx] = alive(x, y, sourceCell, targetCell)
			} else {
				target.s[idx] = dead(x, y, sourceCell, targetCell)
			}
		}
	}

	return nil
}
