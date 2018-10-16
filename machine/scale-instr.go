package machine

import (
	"github.com/moccalotto/nick/field/modifiers"
)

func init() {
	InstructionHandlers["scale"] = Scale
}

func Scale(m *Machine) {
	m.Assert(m.Field != nil, "Cannot evolve a non-initialized field!")
	x := m.ArgAsFloat(0)
	y := x
	if m.ArgCount() > 1 {
		y = m.ArgAsFloat(2)
		m.Assert(m.ArgAsString(1) == "x", "Second arg to scale must be an 'x', but '%s' was given")

	}
	m.Field.Apply(modifiers.NewScaleXY(x, y))
}
