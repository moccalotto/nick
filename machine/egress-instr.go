package machine

import "github.com/moccalotto/nick/field/modifiers"

func init() {
	InstructionHandlers["egress"] = Egress
}

func Egress(m *Machine) {
	m.Assert(m.Field != nil, "Cannot snow a non-initialized field!")

	var thickness, length int
	var direction modifiers.Direction

	switch m.ArgAsString(0) {
	case "north":
		direction = modifiers.North
	case "south":
		direction = modifiers.South
	case "east":
		direction = modifiers.East
	case "west":
		direction = modifiers.West
	case "north-east":
		direction = modifiers.NorthEast
	case "north-west":
		direction = modifiers.NorthWest
	case "south-east":
		direction = modifiers.SouthEast
	case "south-west":
		direction = modifiers.SouthWest
	default:
		m.Assert(
			false,
			"Invalid direction '%s'. Must be one of [north, south, east, west, north-east, north-west, south-east, south-west]",
			m.ArgAsString(0),
		)
	}

	errStr := "Invalid use of egress. Use one of: 'egress [direction]' or 'egress [direction] @ [length] x [thickness]'"

	m.Assert(m.ArgCount() == 1 || m.ArgCount() == 5, errStr)

	if m.ArgCount() == 5 {
		length = m.ArgAsInt(2)
		thickness = m.ArgAsInt(4)
		m.Assert(m.ArgAsString(1) == "@", errStr)
		m.Assert(m.ArgAsString(3) == "x", errStr)
	} else {
		// the default length if the egress is length / 4 squares
		length := min(m.Field.Width(), m.Field.Height()) / 4
		thickness := min(m.Field.Width(), m.Field.Height()) / 10
		if thickness == 0 {
			thickness = 1
		}
		if length == 0 {
			length = 1
		}
	}

	egress := modifiers.NewEgress(direction, length)
	egress.Thickness = thickness

	m.Field.Apply(egress)
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
