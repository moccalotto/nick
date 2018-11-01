package field

type Cell uint8

type CellMapper func(f *Field, x, y int, c Cell) Cell

type CellWalker func(x, y int, c Cell)

func (c Cell) On() bool {
	return c > 0
}

func (c Cell) Toggled() Cell {
	return c ^ 1
}

func (c Cell) AsInt() int {
	return int(c)
}

const OffCell Cell = 0
const LivingCell Cell = 1
