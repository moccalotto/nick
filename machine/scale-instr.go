package machine

import (
	"github.com/moccaloto/nick/field/modifiers"
)

func init() {
	InstructionHandlers["scale"] = Scale
}

func Scale(m *Machine) {
	m.Assert(m.Field != nil, "Cannot evolve a non-initialized field!")

	x := m.ArgAsFloat(0)
	y := x
	if m.HasArg(1) {
		y = m.ArgAsFloat(1)
	}
	m.Field.Apply(modifiers.NewScaleXY(x,y))
}
