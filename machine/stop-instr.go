package machine

import "github.com/moccalotto/nick/field"

func init() {
	InstructionHandlers["stop!"] = Stop
}

func Stop(m *Machine) {
	if m.Cave == nil {
		m.Cave = field.NewField(1, 1)
	}
	m.State.PC = len(m.Tape)
}
