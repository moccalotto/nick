package machine

import "github.com/moccalotto/nick/effects"

func init() {
	InstructionHandlers["rect"] = Rect
}

func Rect(m *Machine) {
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
	r := effects.NewRect(
		m.ArgAsInt(0),
		m.ArgAsInt(1),
		m.ArgAsInt(2),
		m.ArgAsInt(3),
		m.Rng,
	)
	r.Coverage = coverage
	r.On = on
	r.ApplyToField(m.Field)
}
