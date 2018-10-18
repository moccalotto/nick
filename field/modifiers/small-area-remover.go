package modifiers

import (
	"github.com/moccalotto/nick/field"
)

type SmallAreaCRemover struct {
	Threshold int // remove rooms with fewer tiles than this
}

type point struct {
	x, y int
}

func (p point) adjecent() []point {
	return []point{
		point{p.x + 1, p.y}, // east
		point{p.x - 1, p.y}, // west
		point{p.x, p.y + 1}, // north
		point{p.x, p.y - 1}, // south
	}
}

type area []point

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
			if inspected.Alive(x, y) {
				continue
			}

			curState := f.Alive(x, y)

			// retrieve all cells in the given area.
			a := getAreaAround(f, point{x, y})

			if len(a) < m.Threshold {
				// the area was too small. Remove it,
				// and mark its cells as inspected
				for _, p := range a {
					inspected.SetAlive(p.x, p.y, true)
					f.SetAlive(p.x, p.y, !curState)
				}
			} else {
				// the area was large enough to keep,
				// but we still mark the area as inspected.
				for _, p := range a {
					inspected.SetAlive(p.x, p.y, true)
				}
			}
		}
	}

}

func getAreaAround(f *field.Field, p point) area {
	queue := []point{p}
	areaType := f.Alive(p.x, p.y)
	inspected := field.NewField(f.Width(), f.Height())

	result := area{}

	for len(queue) > 0 {
		_p := queue[0]
		queue = queue[1:]

		// anything on the queue can be appended.
		result = append(result, _p)
		inspected.SetAlive(_p.x, _p.y, true)

		for _, c := range _p.adjecent() {
			// outside the map?
			if !f.CoordsInRange(c.x, c.y) {
				continue
			}

			// already inspected?
			if inspected.Alive(c.x, c.y) {
				continue
			}

			// does this cell belong to another area?
			if f.Alive(c.x, c.y) != areaType {
				continue
			}

			// Point has not yet been looked at (or marked for inspection)
			queue = append(queue, c)           // Add c to the queue
			inspected.SetAlive(c.x, c.y, true) // Mark c as inspected.
		}
	}

	return result
}
