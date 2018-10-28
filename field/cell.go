package field

import (
	"image/color"
)

type Cell color.RGBA

type CellMapper func(f *Field, x, y int, c Cell) Cell

type CellWalker func(x, y int, c Cell)

func (c Cell) Alive() bool {
	// Cells with value 1 or greater are considered to be alive
	return c.A == 0xff
}

func (c Cell) Dead() bool {
	// Cells with value 0 or lower are considered to be dead.
	// This means that there are several "dead" states.
	// It also means that you cannot use fields to denote negative heights.
	return c.A == 0x00
}

var DeadCell Cell = Cell{R: 0x00, G: 0x00, B: 0x00, A: 0x00}
var LivingCell Cell = Cell{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
