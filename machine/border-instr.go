package machine

import "github.com/moccalotto/nick/field/modifiers"

func init() {
	InstructionHandlers["border"] = Border
}

func Border(m *Machine) {
	m.Assert(m.Field != nil, "Cannot snow a non-initialized field!")

	border := modifiers.NewBorderSnow(1.0, m.Rng)
	border.Thickness = m.ArgAsInt(0)
	m.Assert(border.Thickness > 0, "Thickness must be > 0")

	if m.HasArg(1) {
		border.Coverage = m.ArgAsFloat(1)
		m.Assert(
			border.Coverage > 0.0 && border.Coverage <= 1.0,
			"Coverage must be between 0% and 100%",
		)
	}

	if m.HasArg(2) {
		m.Assert(
			m.ArgAsString(2) == "(dead)",
			"Only allowed value to this instruction is the string '(dead)'",
		)

		border.Alive = false
	}

	m.Field.Apply(border)
}
