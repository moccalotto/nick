package machine

import (
	"github.com/moccalotto/nick/effects"
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

	t := effects.NewThile(a, b)

	t.ApplyToField(m.Field)
}
