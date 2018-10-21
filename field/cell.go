package field

type Cell int

type CellMapper func(f *Field, x, y int, c Cell) Cell

func (c Cell) Alive() bool {
	// Cells with value 1 or greater are considered to be alive
	return c > 0
}

func (c Cell) Dead() bool {
	// Cells with value 0 or lower are considered to be dead.
	// This means that there are several "dead" states.
	// It also means that you cannot use fields to denote negative heights.
	return c <= 0
}

const LivingCell Cell = 1
const DeadCell Cell = 2
