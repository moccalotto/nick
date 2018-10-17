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

	var depth, width int

	errStr := "Invalid use of egress. Use one of: 'egress [direction]' or 'egress [direction] [width] x [depth]'"

	m.Assert(m.ArgCount() == 1 || m.ArgCount() == 4, errStr)

	if m.ArgCount() == 4 {
		width = m.ArgAsInt(1)
		depth = m.ArgAsInt(3)
		m.Assert(m.ArgAsString(2) == "x", errStr)
	} else {
		// the default width of the egress is width / 4 squares
		width := min(m.Field.Width(), m.Field.Height()) / 4
		depth := min(m.Field.Width(), m.Field.Height()) / 10
		if depth == 0 {
			depth = 1
		}
		if width == 0 {
			width = 1
		}
	}

	direction := makeDirection(m)

	egress := modifiers.NewEgress(direction, width, m.Rng)
	egress.Depth = depth

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
