package effects

import "github.com/moccalotto/nick/field"

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

// an area is all the cells within it.
type area []point

// a room is an area consisting only of the edges-cells of an area.
type room area

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
