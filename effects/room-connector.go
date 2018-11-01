package effects

import (
	"github.com/moccalotto/nick/field"
)

type RoomConnector struct {
	TunnelRadius  float64
	MaxRooms      int
	MaxIterations int
}

func NewRoomConnector(tunnelRadius float64, maxRooms, maxIterations int) *RoomConnector {
	return &RoomConnector{tunnelRadius, maxRooms, maxIterations}
}

func (rc *RoomConnector) getAllRooms(f *field.Field) []field.Area {
	w, h := f.Width(), f.Height()

	result := []field.Area{}
	rawCells := f.Cells()

	// buffer to ensure that we don't look at the same area twice
	inspected := make([]bool, w*h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// On cells are walls
			if rawCells[x+y*w] != field.OffCell {
				continue
			}

			// Don't look at a cell twice
			if inspected[x+y*w] {
				continue
			}

			_area, _ := f.GetAreaAround(x, y)

			r := field.Area{}

			// For each point in the area, check if the given cell is on the edge.
			for _, p := range _area {
				// Mark all cells in the room as inspected
				inspected[p.X+p.Y*w] = true

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

					if rawCells[ap.X+ap.Y*w] == field.OffCell {
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

func (rc *RoomConnector) findClosestPoints(r1, r2 field.Area) (field.Point, field.Point) {
	bestDistSq := uint64(0xFFFFFFFFFFFF)
	var bestP1 field.Point
	var bestP2 field.Point
	for _, p1 := range r1 {
		for _, p2 := range r2 {
			dx := p1.X - p2.X
			dy := p1.Y - p2.Y
			distSq := uint64(dx*dx + dy*dy)

			if distSq < bestDistSq {
				bestDistSq = distSq
				bestP1 = p1
				bestP2 = p2
			}
		}
	}

	return bestP1, bestP2
}

func (rc *RoomConnector) startTunnel(f *field.Field, r1, r2 field.Area) {
	p1, p2 := r1.FindClosestPoints(r2)

	f.SetRadius(p1.X, p1.Y, rc.TunnelRadius, field.OffCell)
	f.SetRadius(p2.X, p2.Y, rc.TunnelRadius, field.OffCell)
}

// Connect all rooms in field
func (rc *RoomConnector) ApplyToField(f *field.Field) {
	for i := 0; i < rc.MaxIterations; i++ {
		roomCount := 0

		rooms := rc.getAllRooms(f)

		// We have reached the allowed number of rooms, exit
		if len(rooms) <= rc.MaxRooms {
			return
		}

		for _, r1 := range rooms {
			roomCount++
			for _, r2 := range rooms {
				if r1[0] == r2[0] {
					continue // don't connect a room to itself.
				}

				rc.startTunnel(f, r1, r2)
			}
		}

	}
}
