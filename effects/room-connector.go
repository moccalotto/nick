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

// Create a tunnel between room1 and room2
// we tunnel from room1 towards room2 and vise versa
// once the two tunnels meet, we stop.
func (rc *RoomConnector) connect(f *field.Field, room1, room2 field.Area) {
	p1, p2, distSq := room1.FindClosestPoints(room2)

	if distSq == 0 {
		return
	}

	area1 := p1.WithinRadius(rc.TunnelRadius)
	area2 := p2.WithinRadius(rc.TunnelRadius)

	f.SetArea(area1, field.OffCell)
	f.SetArea(area2, field.OffCell)

	rc.connect(f, area1, area2)
}

// Connect all rooms in field
func (rc *RoomConnector) ApplyToField(f *field.Field) {
	for i := 0; i < rc.MaxIterations; i++ {
		roomCount := 0

		rooms := f.GetAllRooms()

		// We have reached the allowed number of rooms, exit
		if len(rooms) <= rc.MaxRooms {
			return
		}

		for _, room1 := range rooms {
			roomCount++
			for _, room2 := range rooms {
				if room1[0] == room2[0] {
					continue // don't connect a room to itself.
				}

				rc.connect(f, room1, room2)
			}
		}

	}
}
