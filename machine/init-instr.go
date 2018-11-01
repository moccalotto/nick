package machine

import "github.com/moccalotto/nick/field"

func init() {
	InstructionHandlers["init"] = Init
}

func Init(m *Machine) {
	m.Assert(m.Cave == nil, "You cannot call 'init' more than once!")

	m.Assert(m.ArgAsString(1) == "x", "Args for 'init' must be [number] x [number]")

	m.Cave = field.NewField(
		m.ArgAsInt(0),
		m.ArgAsInt(2),
	)
}
