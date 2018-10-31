package field

type Cell bool

type CellMapper func(f *Field, x, y int, c Cell) Cell

type CellWalker func(x, y int, c Cell)

func (c Cell) On() bool {
	return c == LivingCell
}

func (c Cell) Off() bool {
	return c == OffCell
}

func (c Cell) Toggled() Cell {
	return !c
}

const OffCell Cell = false
const LivingCell Cell = true
