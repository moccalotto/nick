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
	inspected := field.NewField(w, h)

	for y := 0; y < h; y = y + 1 {
		for x := 0; x < w; x = x + 1 {
			// ensure we don't look at the same cell twice
			if a, err := inspected.On(x, y); err != nil {
				panic(err)
			} else if a {
				continue
			}

			curState, _ := f.On(x, y)

			// retrieve all cells in the given area.
			a, err := f.GetAreaAround(x, y)

			if err != nil {
				panic(err)
			}

			if len(a) < m.Threshold {
				// the area was too small. Remove it,
				// and mark its cells as inspected
				for _, p := range a {
					_ = inspected.SetOn(p.X, p.Y, true)
					_ = f.SetOn(p.X, p.Y, !curState)
				}
			} else {
				// the area was large enough to keep,
				// but we still mark the area as inspected.
				for _, p := range a {
					_ = inspected.SetOn(p.X, p.Y, true)
				}
			}
		}
	}
}
