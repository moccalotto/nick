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

	errStr := "Invalid use of egress. Use one of: 'egress [direction]' or 'egress [direction] [size]'"

	m.Assert(m.ArgCount() == 1 || m.ArgCount() == 2, errStr)

	radius := 1.0
	if m.ArgCount() == 2 {
		radius = m.ArgAsFloat(1)
	} else {
		// the default width of the egress is width / 4 squares
		radius = minAsFloat(m.Field.Width(), m.Field.Height()) / 5
	}

	direction := makeDirection(m)

	egress := effects.NewEgress(direction, radius, m.Rng)

	egress.ApplyToField(m.Field)
}

func minAsFloat(a, b int) float64 {
	if a < b {
		return float64(a)
	}

	return float64(b)
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
