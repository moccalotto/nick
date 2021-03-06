package effects

import (
	"github.com/moccalotto/nick/field"
)

type SmallAreaCRemover struct {
	Threshold int // remove rooms with fewer tiles than this
}

func NewSmallAreaCRemover(threshold int) *SmallAreaCRemover {
	return &SmallAreaCRemover{threshold}
}

func (m *SmallAreaCRemover) ApplyToField(f *field.Field) {

	if m.Threshold <= 0 {
		return
	}

	w, h := f.Width(), f.Height()

	// buffer to keep track of all the fields we've looked at
	inspected := make([]bool, w*h)
	rawCells := f.Cells()

	for y := 0; y < h; y = y + 1 {
		for x := 0; x < w; x = x + 1 {
			// ensure we don't look at the same cell twice
			if inspected[x+y*w] {
				continue
			}

			toggled := rawCells[x+y*w].Toggled()

			// retrieve all cells in the given area.
			a, err := f.AreaAround(x, y)

			if err != nil {
				panic(err)
			}

			if len(a) < m.Threshold {
				// the area was too small. Remove it,
				// and mark its cells as inspected
				for _, p := range a {
					inspected[p.X+p.Y*w] = true
					rawCells[p.X+p.Y*w] = toggled
				}
			} else {
				// the area was large enough to keep,
				// but we still mark the area as inspected.
				for _, p := range a {
					inspected[p.X+p.Y*w] = true
				}
			}
		}
	}
}
