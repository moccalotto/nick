package field

type Cell bool

type CellMapper func(f *Field, x, y int, c Cell) Cell

type CellWalker func(x, y int, c Cell)

func (c Cell) Alive() bool {
	return c == LivingCell
}

func (c Cell) Dead() bool {
	return c == DeadCell
}

const DeadCell Cell = false
const LivingCell Cell = true
