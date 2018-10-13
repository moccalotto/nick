package machine

import "github.com/moccalotto/nick/field/modifiers"

func init() {
	InstructionHandlers["border"] = Border
}

func Border(m *Machine) {
	m.Assert(m.Field != nil, "Cannot snow a non-initialized field!")

	border := modifiers.NewBorderSnow(1.0)
	border.Thickness = m.ArgAsInt(0)
	m.Assert(border.Thickness > 0, "Thickness must be > 0")

	if m.HasArg(2) {
		m.Assert(m.ArgAsString(1) == "@", "Second arg to 'border' must be an '@', but '%s' was given")
		border.Coverage = m.ArgAsFloat(2)
		m.Assert(
			border.Coverage > 0.0 && border.Coverage <= 1.0,
			"Coverage must be between 0% and 100%",
		)
	}

	m.Field.Apply(border)
}
