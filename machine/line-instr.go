package machine

import "github.com/moccalotto/nick/effects"
import "github.com/moccalotto/nick/field"

func init() {
	InstructionHandlers["line"] = Line
}

// Line:
// usage: line x0 y0 x1 y1
func Line(m *Machine) {
	m.Assert(m.Cave != nil, "The '%s' instruction can only be used after the 'init'", m.CurrentInstruction().Cmd)

	if m.HasArg(4) {
		m.ArgAsFloat(4)
	}

	l := effects.NewLine(
		m.ArgAsInt(0),
		m.ArgAsInt(1),
		m.ArgAsInt(2),
		m.ArgAsInt(3),
		m.Rng,
	)
	if m.HasArg(5) {
		if m.ArgAsString(5) == "(off)" {
			l.Cell = field.OffCell
		} else {
			m.Throw("The only allowed value for the fifth argument is the string '(off)'")
		}

	}

	l.ApplyToField(m.Cave)
}
