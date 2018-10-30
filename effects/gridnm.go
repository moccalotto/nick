package effects

import (
	"github.com/moccalotto/nick/field"
	"math"
)

// GridNM represents a grid that
// has a fixed number of rows and columns,
// no matter how big the field is.
type GridNM struct {
	Columns int
	Rows    int
}

// Create a new grid with rows-1 horizontal lines and cols-1 vertical lines.
func NewGridNM(cols, rows int) *GridNM {
	return &GridNM{
		Columns: cols,
		Rows:    rows,
	}
}

// When dividing a line into n+1 parts using n lines,
// we calculate the positions (one-dimensional coordinates)
// of each line. We return it as a lookup-table. I.e. a
// map of coordinate => true
func makePositions(length, segments int) map[int]bool {
	delta := float64(length) / float64(segments)

	var x float64

	result := map[int]bool{}

	lineCount := segments - 1

	for i := 0; i < lineCount; i++ {
		x += delta

		result[int(math.Floor(x))] = true
	}

	return result
}

func (grid *GridNM) ApplyToField(f *field.Field) {
	xs := makePositions(f.Width(), grid.Columns)
	ys := makePositions(f.Height(), grid.Rows)

	f.Map(func(f *field.Field, x, y int, c field.Cell) field.Cell {

		// we don't draw grids on top of living cells.
		if c.On() {
			return c
		}

		if v := xs[x]; v {
			return field.LivingCell
		}

		if v := ys[y]; v {
			return field.LivingCell
		}

		return c
	})
}
