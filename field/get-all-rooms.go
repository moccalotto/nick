package field

func (f *Field) GetAllRooms() []Area {
	result := []Area{}

	// buffer to ensure that we don't look at the same area twice
	inspected := make([]bool, len(f.s))

	for y := 0; y < f.h; y++ {
		for x := 0; x < f.w; x++ {
			// Only off-cells can be considered to be rooms
			// All other cells are inaccessible.
			if f.s[x+y*f.w] != OffCell {
				continue
			}

			// Don't look at same cell twice
			if inspected[x+y*f.w] {
				continue
			}

			area, _ := f.GetAreaAround(x, y)

			r := Area{}

			// For each point in the area, check if the given cell is on the edge.
			for _, p := range area {
				// Mark all cells in the room as inspected
				inspected[p.X+p.Y*f.w] = true

				// check if all adjacent cells are also in the room
				// if not, then it's a point on the edge, and therefore
				// constitunes the room
				for _, ap := range p.Adjacent() {
					if !f.CoordsInRange(ap.X, ap.Y) {
						// note: this means that cells on the very edge of a field
						// will not be marked as edges.
						continue
					}

					// If at least one adjacent cell is off (i.e. outside the room)
					// The current cell must be on the edge of the room

					if f.s[ap.X+ap.Y*f.w] == OffCell {
						r = append(r, p)
						continue
					}
				}
			}

			if len(r) > 0 {
				result = append(result, r)
			}
		}
	}

	return result
}
