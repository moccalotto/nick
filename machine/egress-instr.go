package machine

import (
	"github.com/moccalotto/nick/effects"
	"sort"
)

var dirs map[string]effects.Direction = map[string]effects.Direction{
	"random":     effects.Random,
	"north":      effects.North,
	"N":          effects.North,
	"north-east": effects.NorthEast,
	"NE":         effects.NorthEast,
	"east":       effects.East,
	"E":          effects.East,
	"south-east": effects.SouthEast,
	"SE":         effects.SouthEast,
	"south":      effects.South,
	"S":          effects.South,
	"south-west": effects.SouthWest,
	"SW":         effects.SouthWest,
	"west":       effects.West,
	"W":          effects.West,
	"north-west": effects.NorthWest,
	"NW":         effects.NorthWest,
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

	egress := effects.NewEgress(direction, width, m.Rng)
	egress.Depth = depth

	m.Field.Apply(egress)
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func makeDirection(m *Machine) effects.Direction {
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
