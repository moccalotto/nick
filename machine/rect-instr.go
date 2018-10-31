package machine

import "github.com/moccalotto/nick/effects"
import "github.com/moccalotto/nick/field"

func init() {
	InstructionHandlers["rect"] = Rect
}

func Rect(m *Machine) {
	r := effects.NewRect(
		m.ArgAsInt(0),
		m.ArgAsInt(1),
		m.ArgAsInt(2),
		m.ArgAsInt(3),
		m.Rng,
	)
	if m.HasArg(4) {
		if m.ArgAsString(4) == "(off)" {
			r.Cell = field.OffCell
		} else {
			m.Throw("The only allowed value for the fifth argument is the string '(off)'")
		}

	}
	r.ApplyToField(m.Field)
}
