package machine

import (
	"github.com/moccalotto/nick/field/modifiers"
)

func init() {
	InstructionHandlers["connect-rooms"] = ConnectRooms
}

//  connect-rooms 3	// connect rooms with tunnels of diameter 3
func ConnectRooms(m *Machine) {
	maxRooms := 1
	maxIterations := 5

	radius := m.ArgAsFloat(0)

	m.Assert(radius > 0, "Radius must be > 0")

	if radius < 1.0 {
		if m.Field.Width() < m.Field.Height() {
			radius = radius * float64(m.Field.Width())
		} else {
			radius = radius * float64(m.Field.Height())
		}
	}

	if m.ArgCount() > 1 {
		maxRooms = m.ArgAsInt(1)
	}
	if m.ArgCount() > 2 {
		maxIterations = m.ArgAsInt(2)
	}

	rc := modifiers.NewRoomConnector(
		radius,
		maxRooms,
		maxIterations,
	)

	rc.ApplyToField(m.Field)
}
