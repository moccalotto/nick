package field

import (
	"fmt"
	"image/color"
	"sync"
)

// Field represents a two-dimensional field of cells.
type Field struct {
	s    []Cell // cells
	w, h int
	// should we allow "circular overflow?"
	Outside Cell // if trying to access a cell outside the area, is it on or off?

	Palette color.Palette
}

// NewField returns an empty field of the specified width and height.
func NewField(w, h int) *Field {
	return &Field{
		s:       make([]Cell, h*w),
		w:       w,
		h:       h,
		Outside: OnCell,
		Palette: DefaultPalette(),
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
func (f *Field) Get(x, y int) (Cell, error) {
	if err := f.errCoordsInRange(x, y); err != nil {
		return f.Outside, err
	}

	return f.s[y*f.w+x], nil
}

// Turn on all cells in the area
func (f *Field) SetRadius(x, y int, r float64, c Cell) {
	area := (Point{x, y}).WithinRadius(r)

	f.SetArea(area, c)
}

// Set multiple points
func (f *Field) SetArea(area Area, c Cell) {
	for _, p := range area {
		_ = f.Set(p.X, p.Y, c)
	}
}

// On reports whether the specified cell is on.
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (f *Field) On(x, y int) (bool, error) {
	if err := f.errCoordsInRange(x, y); err != nil {
		return f.Outside.On(), err
	}

	return f.s[y*f.w+x].On(), nil
}

func (f *Field) Width() int {
	return f.w
}

func (f *Field) Height() int {
	return f.h
}

func (f *Field) AspectRatio() float64 {
	return float64(f.w) / float64(f.h)
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
	_y := y - 1
	for _x := x - 1; _x <= x+1; _x++ {
		c, _ := f.Get(_x, _y)
		neighbourCount += c.AsInt()
	}

	// Check neighbours on the line below
	_y = y + 1
	for _x := x - 1; _x <= x+1; _x++ {
		c, _ := f.Get(_x, _y)
		neighbourCount += c.AsInt()
	}

	// Check neighbour to the left
	c, _ := f.Get(x-1, y)
	neighbourCount += c.AsInt()

	// Check neighbour to the right
	c, _ = f.Get(x+1, y)
	neighbourCount += c.AsInt()

	return neighbourCount
}

// Call a function on each cell
func (f *Field) Walk(w CellWalker) {
	for y := 0; y < f.h; y++ {
		for x := 0; x < f.w; x++ {
			idx := y*f.w + x
			w(x, y, f.s[idx])
		}
	}
}

// Call a function on each cell (async)
func (f *Field) WalkAsync(w CellWalker) {
	var wg sync.WaitGroup
	for y := 0; y < f.h; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := 0; x < f.Width(); x++ {
				idx := y*f.w + x
				w(x, y, f.s[idx])
			}
		}(y)
	}

	wg.Wait()
}

// Map each cell to another value
func (f *Field) Map(m CellMapper) {
	s := make([]Cell, len(f.s))

	for y := 0; y < f.h; y++ {
		for x := 0; x < f.w; x++ {
			idx := y*f.w + x
			s[idx] = m(
				f,
				x,
				y,
				f.s[idx],
			)
		}
	}

	f.s = s
}

// Map each cell to another value, but do it asynchornously
func (f *Field) MapAsync(m CellMapper) {
	f.s = f.MapAsyncToNewField(m).s
}

// Map each cell to a cell in a new field with same properties.
func (f *Field) MapAsyncToNewField(m CellMapper) *Field {
	s := make([]Cell, len(f.s))

	var wg sync.WaitGroup
	for y := 0; y < f.h; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := 0; x < f.Width(); x++ {
				idx := y*f.w + x
				s[idx] = m(
					f,
					x,
					y,
					f.s[idx],
				)
			}
		}(y)
	}

	wg.Wait()

	return &Field{
		s:       s,
		w:       f.w,
		h:       f.h,
		Outside: f.Outside,
		Palette: f.Palette,
	}
}

func (f *Field) HasEgress() bool {
	rawCells := f.Cells()
	h, w := f.Height(), f.Width()

	// check if north face has 1 or more dead cells
	for x := 0; x < w; x++ {
		if rawCells[x] == OffCell {
			return true
		}
	}

	// check if east face has 1 or more dead cells
	x := w - 1
	for y := 0; y < h; y++ {
		if rawCells[x+y*w] == OffCell {
			return true
		}
	}

	// check if south face has 1 or more dead cells
	y := h - 1
	for x := 0; x < w; x++ {
		if rawCells[x+y*w] == OffCell {
			return true
		}
	}

	// check if west face has 1 or more dead cells
	for y := 0; y < h; y++ {
		if rawCells[y*w] == OffCell {
			return true
		}
	}

	return false
}
