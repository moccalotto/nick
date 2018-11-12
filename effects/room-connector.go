package effects

import (
	"github.com/moccalotto/nick/field"
)

type RoomConnector struct {
	TunnelRadius float64
	MaxRooms     int
}

func NewRoomConnector(tunnelRadius float64, maxRooms, maxIterations int) *RoomConnector {
	if maxRooms < 1 {
		panic("maxRooms must be ≥ 1")
	}
	return &RoomConnector{tunnelRadius, maxRooms}
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
// Algorithm:
// let MaxRooms = the max number of rooms allowed in the output
// let roomCount = the number of rooms found
// if roomCount <= MaxRooms
//     stop
// else
//     create tunnel between two arbitrary rooms

func (rc *RoomConnector) ApplyToField(f *field.Field) {

	for {
		// we need to fetch all rooms on every iteration
		// because tunnelling from one room to another
		// may have connected other rooms as well.
		rooms := f.AllRoomWalls()
		roomCount := len(rooms)

		if roomCount <= rc.MaxRooms {
			return
		}

		lastRoom := rooms[roomCount-1]
		closestRoom, _, _, _ := lastRoom.FindClosestArea(rooms[0 : roomCount-1])

		rc.connect(f, lastRoom, closestRoom)

		// semantically, this statement has no effect,
		// it is only there for optimization purposes.
		if roomCount == rc.MaxRooms+1 {
			return
		}
	}

}
