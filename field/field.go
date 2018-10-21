package field

import "fmt"

type Cell int

type Point struct {
	x int
	y int
}

type CellHandler func(x, y int, c Cell) Cell

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

func LivingCell() Cell {
	return Cell(1)
}

func DeadCell() Cell {
	return Cell(0)
}

// Field represents a two-dimensional field of cells.
type Field struct {
	s       []Cell // cells
	w, h    int
	outside Cell // are the cells outside the scope defined by [0..w]x[0..y] alive ?
}

type Modifier interface {
	ApplyToField(f *Field)
}

// NewField returns an empty field of the specified width and height.
func NewField(w, h int) *Field {
	s := make([]Cell, h*w)
	return &Field{s, w, h, LivingCell()}
}

func (f *Field) Apply(m Modifier) {
	m.ApplyToField(f)
}

func (f *Field) CoordsInRange(x, y int) bool {
	return x < f.w &&
		y < f.h &&
		x >= 0 &&
		y >= 0
}

func (f *Field) CoordsMustBeInRange(x, y int) {
	if !f.CoordsInRange(x, y) {
		panic(fmt.Sprintf(
			"Coords [%d, %d] are out of range [0..%d, 0..%d]",
			x,
			y,
			f.w-1,
			f.h-1,
		))
	}
}

// Set sets the state of the specified cell to the given value.
func (f *Field) Set(x, y int, c Cell) {
	f.CoordsMustBeInRange(x, y)

	f.s[y*f.w+x] = c
}

// Set sets the state of the specified cell to the given value.
func (f *Field) SetAlive(x, y int, b bool) {
	if b {
		f.Set(x, y, LivingCell())
	} else {
		f.Set(x, y, DeadCell())
	}

}

func (f *Field) SetAliveRadius(x, y int, r float64, b bool) {

	r2 := r * r
	for i := 0; float64(i) <= r; i++ {
		for j := 0; j <= i; j++ {
			d2 := i*i + j*j

			// point outside radius
			if float64(d2) > r2 {
				continue
			}

			f.robustSetAlive(x, y, b)

			f.robustSetAlive(x+i, y+j, b)
			f.robustSetAlive(x+i, y-j, b)
			f.robustSetAlive(x-i, y+j, b)
			f.robustSetAlive(x-i, y-j, b)

			f.robustSetAlive(x+j, y+i, b)
			f.robustSetAlive(x+j, y-i, b)
			f.robustSetAlive(x-j, y+i, b)
			f.robustSetAlive(x-j, y-i, b)
		}
	}
}

// Alive reports whether the specified cell is alive.
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (f *Field) Alive(x, y int) bool {
	f.CoordsMustBeInRange(x, y)
	return f.s[y*f.w+x].Alive()
}

func (f *Field) robustGet(x, y int) Cell {
	if !f.CoordsInRange(x, y) {
		return f.outside
	}

	return f.s[y*f.w+x]
}

func (f *Field) robustSetAlive(x, y int, b bool) {
	if !f.CoordsInRange(x, y) {
		return
	}

	f.SetAlive(x, y, b)
}

func (f *Field) robustAlive(x, y int) bool {
	return f.robustGet(x, y).Alive()
}

func (f *Field) Dead(x, y int) bool {
	return !f.Alive(x, y)
}

func (f *Field) Width() int {
	return f.w
}

func (f *Field) Height() int {
	return f.h
}

func (f *Field) Cells() []Cell {
	return f.s
}

func (f *Field) SetCells(w, h int, s []Cell) {
	if len(s) != w*h {
		panic("Invalid use of SetCells(). w*h must be len(s)")
	}
	f.s = s
	f.h = h
	f.w = w
}

func (f *Field) NeighbourCount(x, y int) int {
	neighbourCount := 0

	// Check neighbours above
	for _x := x - 1; _x <= x+1; _x++ {
		if f.robustAlive(_x, y-1) {
			neighbourCount++
		}
	}
	// Check neighbours on the line below
	for _x := x - 1; _x <= x+1; _x++ {
		if f.robustAlive(_x, y+1) {
			neighbourCount++
		}
	}

	// Check neighbour to the left
	if f.robustAlive(x-1, y) {
		neighbourCount++
	}

	// Check neighbourCount to the right
	if f.robustAlive(x+1, y) {
		neighbourCount++
	}

	return neighbourCount
}
