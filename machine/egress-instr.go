package machine

import (
	"github.com/moccalotto/nick/field/modifiers"
	"sort"
)

var dirs map[string]modifiers.Direction = map[string]modifiers.Direction{
	"random":     modifiers.Random,
	"north":      modifiers.North,
	"N":          modifiers.North,
	"north-east": modifiers.NorthEast,
	"NE":         modifiers.NorthEast,
	"east":       modifiers.East,
	"E":          modifiers.East,
	"south-east": modifiers.SouthEast,
	"SE":         modifiers.SouthEast,
	"south":      modifiers.South,
	"S":          modifiers.South,
	"south-west": modifiers.SouthWest,
	"SW":         modifiers.SouthWest,
	"west":       modifiers.West,
	"W":          modifiers.West,
	"north-west": modifiers.NorthWest,
	"NW":         modifiers.NorthWest,
}

func init() {
	InstructionHandlers["egress"] = Egress
}

func Egress(m *Machine) {
	m.Assert(m.Field != nil, "Cannot snow a non-initialized field!")

	var thickness, length int

	errStr := "Invalid use of egress. Use one of: 'egress [direction]' or 'egress [direction] [length] x [thickness]'"

	m.Assert(m.ArgCount() == 1 || m.ArgCount() == 4, errStr)

	if m.ArgCount() == 4 {
		length = m.ArgAsInt(1)
		thickness = m.ArgAsInt(3)
		m.Assert(m.ArgAsString(2) == "x", errStr)
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

	direction := makeDirection(m)

	egress := modifiers.NewEgress(direction, length, m.Rng)
	egress.Thickness = thickness

	m.Field.Apply(egress)
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func makeDirection(m *Machine) modifiers.Direction {
	direction, ok := dirs[m.ArgAsString(0)]

	if !ok {
		keys := make([]string, 0, len(dirs))
		for k := range dirs {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		m.Assert(
			false,
			"Invalid direction '%s'. Must be one of %v",
			m.ArgAsString(0),
			keys,
		)
	}

	return direction
}
