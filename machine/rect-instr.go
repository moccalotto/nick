package machine

import "github.com/moccalotto/nick/field/modifiers"

func init() {
	InstructionHandlers["rect"] = Rect
}

func Rect(m *Machine) {
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
	r := modifiers.NewRect(
		m.ArgAsInt(0),
		m.ArgAsInt(1),
		m.ArgAsInt(2),
		m.ArgAsInt(3),
		m.Rng,
	)
	r.Coverage = coverage
	r.Alive = alive
	r.ApplyToField(m.Field)
}
