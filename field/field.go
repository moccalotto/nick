package field

// Field represents a two-dimensional field of cells.
type Field struct {
	s       [][]bool // cells
	w, h    int
	outside bool // are the cells outside the scope defined by [0..w]x[0..y] alive ?
}

type Modifier interface {
	ApplyToField(f *Field)
}

// NewField returns an empty field of the specified width and height.
func NewField(w, h int) *Field {
	s := make([][]bool, h)
	for i := range s {
		s[i] = make([]bool, w)
	}
	return &Field{s, w, h, true}
}

func (f *Field) Apply(m Modifier) {
	m.ApplyToField(f)
}

// Set sets the state of the specified cell to the given value.
func (f *Field) Set(x, y int, b bool) {
	f.s[y][x] = b
	// f.s[(y+f.h)%f.h][(x + f.h%f.h)] = b
}

// Alive reports whether the specified cell is alive.
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (f *Field) Alive(x, y int) bool {
	return f.s[y][x]
}

func (f *Field) robustAlive(x, y int) bool {
	if x >= f.w {
		return f.outside
	}
	if x < 0 {
		return f.outside
	}
	if y >= f.h {
		return f.outside
	}
	if y < 0 {
		return f.outside
	}
	return f.Alive(x, y)
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

func (f *Field) Cells() [][]bool {
	return f.s
}

func (f *Field) SetCells(s [][]bool) {
	f.s = s
	f.h = len(s)
	f.w = len(s[0])
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
