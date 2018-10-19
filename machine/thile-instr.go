package machine

import (
	"github.com/moccalotto/nick/field/modifiers"
)

func init() {
	InstructionHandlers["thile"] = Thile
}

func Thile(m *Machine) {
	a := m.ArgAsInt(0)
	b := 0
	if m.HasArg(1) {
		b = m.ArgAsInt(1)
	}

	t := modifiers.NewThile(a, b)

	t.ApplyToField(m.Field)
}
