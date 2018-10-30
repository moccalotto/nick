package machine

import "github.com/moccalotto/nick/effects"

func init() {
	InstructionHandlers["line"] = Line
}

// Line:
// usage: line x0 y0 x1 y1
func Line(m *Machine) {
	coverage := 1.0

	if m.HasArg(4) {
		m.ArgAsFloat(4)
	}

	on := true

	if m.HasArg(5) {
		if m.ArgAsString(5) == "(off)" {
			on = false
		} else {
			m.Throw("The only allowed value for the fifth argument is the string '(off)'")
		}

	}
	l := effects.NewLine(
		m.ArgAsInt(0),
		m.ArgAsInt(1),
		m.ArgAsInt(2),
		m.ArgAsInt(3),
		m.Rng,
	)
	l.On = on
	l.Coverage = coverage

	l.ApplyToField(m.Field)
}
