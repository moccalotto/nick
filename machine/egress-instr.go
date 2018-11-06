package machine

import (
	"github.com/moccalotto/nick/effects"
	"sort"
)

var dirs map[string]effects.Direction = map[string]effects.Direction{
	"random": effects.Random,
	"north":  effects.North,
	"N":      effects.North,
	"east":   effects.East,
	"E":      effects.East,
	"south":  effects.South,
	"S":      effects.South,
	"west":   effects.West,
	"W":      effects.West,
}

func init() {
	InstructionHandlers["egress"] = Egress
	InstructionHandlers["ensure-egress"] = EnsureEgress
}

func EnsureEgress(m *Machine) {
	if m.Cave.HasEgress() {
		return
	}

	effects.NewEgress(
		effects.Random,
		minAsFloat(m.Cave.Width(), m.Cave.Height())*0.1,
		m.Rng,
	).ApplyToField(m.Cave)
}

func Egress(m *Machine) {
	m.Assert(m.Cave != nil, "The '%s' instruction can only be used after the 'init'", m.CurrentInstruction().Cmd)

	errStr := "Invalid use of egress. Use one of: 'egress [direction]' or 'egress [direction] [size]'"

	m.Assert(m.ArgCount() == 1 || m.ArgCount() == 2, errStr)

	radius := 1.0
	if m.ArgCount() == 2 {
		radius = m.ArgAsFloat(1)
	} else {
		radius = 1.0 / 6.0
	}

	// the radius was specified in the range (0, 1)
	// set the radius as a factor of the smalles cave dimension
	if radius < 1 {
		radius = radius * minAsFloat(m.Cave.Width(), m.Cave.Height()) / 2.0
	}

	if radius < 1 {
		radius = 1
	}

	direction := makeDirection(m)

	egress := effects.NewEgress(direction, radius, m.Rng)

	egress.ApplyToField(m.Cave)
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
