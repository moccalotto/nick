package machine

import (
	"github.com/moccalotto/nick/effects"
)

func init() {
	InstructionHandlers["connect-rooms"] = ConnectRooms
}

// connect-rooms 3
// connect-rooms 3 20
func ConnectRooms(m *Machine) {
	radius := m.ArgAsFloat(0)

	m.Assert(radius > 0, "Radius must be > 0")

	if radius < 1.0 {
		if m.Cave.Width() > m.Cave.Height() {
			radius = radius * float64(m.Cave.Height())
		} else {
			radius = radius * float64(m.Cave.Width())
		}
	}

	rc := effects.NewRoomConnector(
		radius,
		1,
		10,
	)

	rc.ApplyToField(m.Cave)
}
