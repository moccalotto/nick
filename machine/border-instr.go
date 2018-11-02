package machine

import (
	"github.com/moccalotto/nick/effects"
	"github.com/moccalotto/nick/field"
)

func init() {
	InstructionHandlers["border"] = Border
}

func Border(m *Machine) {
	m.Assert(m.Cave != nil, "The '%s' instruction can only be used after the 'init'", m.CurrentInstruction().Cmd)

	border := effects.NewBorderSnow(1.0, m.Rng)

	thickness := m.ArgAsFloat(0)
	m.Assert(thickness > 0, "Thickness must be > 0")

	// if the thickness is a float in the range (0, 1), it will be converted to
	// a percentage of the smallest diameter of the cave. So if the thickness is
	// .10, then the border size (in cells) will be 10% of the lower of the width or height,
	if thickness < 1.0 && m.Cave.Width() > m.Cave.Height() {
		border.Thickness = int(thickness * float64(m.Cave.Height()))
	} else if thickness < 1.0 {
		border.Thickness = int(thickness * float64(m.Cave.Width()))
	} else {
		border.Thickness = int(thickness)
	}

	if m.HasArg(1) {
		border.Coverage = m.ArgAsFloat(1)
		m.Assert(
			border.Coverage > 0.0 && border.Coverage <= 1.0,
			"Coverage must be between 0% and 100%",
		)
	}

	if m.HasArg(2) {
		m.Assert(
			m.ArgAsString(2) == "(off)",
			"Only allowed value to this instruction is the string '(off)'",
		)

		border.Cell = field.OffCell
	}

	border.ApplyToField(m.Cave)
}
