package machine

import (
	"github.com/moccalotto/nick/effects"
)

func init() {
	InstructionHandlers["thiele"] = Thiele
}

func Thiele(m *Machine) {
	a := m.ArgAsInt(0)
	b := 0
	if m.HasArg(1) {
		b = m.ArgAsInt(1)
	}

	m.Assert(a > 0, "First argument must be an integer > 0")
	m.Assert(b >= 0, "Second argument (if present) must be an integer ≥ 0")

	effects.NewThielePattern(a, b).ApplyToField(m.Cave)
}
