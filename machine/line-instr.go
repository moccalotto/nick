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

	alive := true

	if m.HasArg(5) {
		if m.ArgAsString(5) == "(dead)" {
			alive = false
		} else {
			m.Throw("The only allowed value for the fifth argument is the string '(dead)'")
		}

	}
	l := effects.NewLine(
		m.ArgAsInt(0),
		m.ArgAsInt(1),
		m.ArgAsInt(2),
		m.ArgAsInt(3),
		m.Rng,
	)
	l.Alive = alive
	l.Coverage = coverage

	l.ApplyToField(m.Field)
}
