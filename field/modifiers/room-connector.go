package modifiers

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

func (rc *RoomConnector) getAllRooms(f *field.Field) []room {
	w, h := f.Width(), f.Height()

	result := []room{}

	// buffer to ensure that we don't look at the same area twice
	inspected := field.NewField(w, h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Alive cells are walls
			if f.Alive(x, y) {
				continue
			}

			// Don't look at a cell twice
			if inspected.Alive(x, y) {
				continue
			}

			_area := getAreaAround(f, point{x, y})

			r := room{}

			// For each point in the area, check if the given cell is on the edge.
			for _, p := range _area {
				// Mark all cells in the room as inspected
				inspected.SetAlive(p.x, p.y, true)

				// check if all adjacent cells are also in the room
				// if not, then it's a point on the edge, and therefore
				// constitunes the room
				for _, ap := range p.adjecent() {
					if !f.CoordsInRange(ap.x, ap.y) {
						// note: this means that cells on the very edge of a field
						// will not be marked as edges.
						continue
					}

					// If at least one adjacent cell is dead (i.e. outside the room)
					// The current cell must be on the edge of the room
					if f.Dead(ap.x, ap.y) {
						r = append(r, p)
						continue
					}
				}
			}

			result = append(result, r)
		}
	}

	return result
}

func (rc *RoomConnector) findClosestPoints(r1, r2 room) (point, point) {
	bestDistSq := uint64(0xFFFFFFFFFFFF)
	var bestP1 point
	var bestP2 point
	for _, p1 := range r1 {
		for _, p2 := range r2 {
			dx := p1.x - p2.x
			dy := p1.y - p2.y
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

func (rc *RoomConnector) startTunnel(f *field.Field, r1, r2 room) {
	p1, p2 := rc.findClosestPoints(r1, r2)

	f.SetAliveRadius(p1.x, p1.y, rc.TunnelRadius, false)
	f.SetAliveRadius(p2.x, p2.y, rc.TunnelRadius, false)
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
