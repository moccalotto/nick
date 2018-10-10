package field

// Field represents a two-dimensional field of cells.
type Field struct {
	s    [][]bool // cells
	w, h int
}

// NewField returns an empty field of the specified width and height.
func NewField(w, h int) *Field {
	s := make([][]bool, h)
	for i := range s {
		s[i] = make([]bool, w)
	}
	return &Field{s, w, h}
}

// Set sets the state of the specified cell to the given value.
func (f *Field) Set(x, y int, b bool) {
	f.s[(y + f.h) % f.h][(x + f.h % f.h)] = b
}

// Alive reports whether the specified cell is alive.
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (f *Field) Alive(x, y int) bool {
	return f.s[(y + f.h) % f.h][(x + f.w) % f.w]
}

// Next returns the state of the specified cell at the next time step.
func (f *Field) NextCellState(x, y int, r *Rules) bool {
	// Count the adjacent cells that are alive.
	neighbourCount := 0

	// TODO: optimize.
	// 	Inspect the entire "line" of 3 cells above the current cell in one go.
	//	Inspect the entire "line" of 3 cells below the current cell in one go.
	//	Inspect left neighbour.
	//	Inspect right neighbour.
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && f.Alive(x+i, y+j) {
				neighbourCount++
			}
		}
	}

	return r.NextCellState(
		f.Alive(x, y),
		neighbourCount,
	)
}

func (f *Field) Evolve(r *Rules) *Field {
	next := NewField(f.w, f.h)

	for x := 0; x < f.w; x++ {
		for y := 0; y < f.w; y++ {
			next.Set(x, y, f.NextCellState(x, y, r))
		}
	}

	return next
}
