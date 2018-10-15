package machine

import "github.com/moccalotto/nick/field/modifiers"
import "sort"

func init() {
	InstructionHandlers["egress"] = Egress
}

func Egress(m *Machine) {
	m.Assert(m.Field != nil, "Cannot snow a non-initialized field!")

	var thickness, length int

	dirs := map[string]modifiers.Direction{
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
