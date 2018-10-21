package machine

import "github.com/moccalotto/nick/effects"

func init() {
	InstructionHandlers["snow"] = Snow
}

func Snow(m *Machine) {
	m.Assert(m.Field != nil, "Cannot snow a non-initialized field!")

	// TODO: allow a "negative" or "dead" modifier to the snow command.

	probability := m.ArgAsFloat(0)

	m.Assert(
		probability >= 0.0 && probability <= 1,
		"Snow takes a number in the range [0, 1] - %.f was given",
		probability,
	)

	snow := effects.NewSnow(probability, m.Rng)

	if m.HasArg(1) {
		arg1 := m.ArgAsString(1)
		m.Assert(
			arg1 == "(dead)",
			"The only value allowed for the optional second argument is the string '(dead)'. The string '%s' was provided",
			arg1,
		)
		snow.Alive = false
	}

	m.Field.Apply(snow)
}
