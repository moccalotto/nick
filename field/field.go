package field

import "fmt"

// Field represents a two-dimensional field of cells.
type Field struct {
	s    []Cell // cells
	w, h int
	// should we allow "circular overflow?"
	outside Cell // if trying to access a cell outside the area, is it alive or dead?
}

// NewField returns an empty field of the specified width and height.
func NewField(w, h int) *Field {
	return &Field{
		s:       make([]Cell, h*w),
		w:       w,
		h:       h,
		outside: LivingCell,
	}
}

// Are the coordinates in range?
func (f *Field) CoordsInRange(x, y int) bool {
	return x < f.w &&
		y < f.h &&
		x >= 0 &&
		y >= 0
}

// Return an error if the coordinates are not in reange.
func (f *Field) errCoordsInRange(x, y int) error {
	if !f.CoordsInRange(x, y) {
		return fmt.Errorf(
			"Coords [%d, %d] are out of range [0..%d, 0..%d]",
			x,
			y,
			f.w-1,
			f.h-1,
		)
	}

	return nil
}

// Panic if coordinates are not in range
func (f *Field) CoordsMustBeInRange(x, y int) {
	if err := f.errCoordsInRange(x, y); err != nil {
		panic(err)
	}
}

// Set sets the state of the specified cell to the given value.
func (f *Field) Set(x, y int, c Cell) error {
	if err := f.errCoordsInRange(x, y); err != nil {
		return err
	}

	f.s[y*f.w+x] = c

	return nil
}

// Set sets the state of the specified cell to the given value.
func (f *Field) SetAlive(x, y int, b bool) error {
	if b {
		return f.Set(x, y, LivingCell)
	} else {
		return f.Set(x, y, DeadCell)
	}

}

// Set all cells in the area to be alive.
func (f *Field) SetAliveRadius(x, y int, r float64, b bool) {
	points := (Point{x, y}).WithinRadius(r)
	cell := DeadCell
	if b {
		cell = LivingCell
	}

	for _, p := range points {
		_ = f.Set(p.X, p.Y, cell)
	}
}

// Alive reports whether the specified cell is alive.
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (f *Field) Alive(x, y int) (bool, error) {
	if err := f.errCoordsInRange(x, y); err != nil {
		return f.outside.Alive(), err
	}

	return f.s[y*f.w+x].Alive(), nil
}

func (f *Field) Dead(x, y int) (bool, error) {
	a, err := f.Alive(x, y)

	return !a, err
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

func (f *Field) ReplaceCells(w, h int, s []Cell) {
	if len(s) != w*h {
		panic("Invalid use of ReplaceCells(). w*h must be len(s)")
	}
	f.s = s
	f.h = h
	f.w = w
}

// Get number of neighbours
func (f *Field) NeighbourCount(x, y int) int {
	neighbourCount := 0

	// Check neighbours above
	for _x := x - 1; _x <= x+1; _x++ {
		if a, _ := f.Alive(_x, y-1); a {
			neighbourCount++
		}
	}
	// Check neighbours on the line below
	for _x := x - 1; _x <= x+1; _x++ {
		if a, _ := f.Alive(_x, y+1); a {
			neighbourCount++
		}
	}

	// Check neighbour to the left
	if a, _ := f.Alive(x-1, y); a {
		neighbourCount++
	}

	// Check neighbourCount to the right
	if a, _ := f.Alive(x+1, y); a {
		neighbourCount++
	}

	return neighbourCount
}
