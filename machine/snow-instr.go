package machine

import "github.com/moccalotto/nick/effects"
import "github.com/moccalotto/nick/field"

func init() {
	InstructionHandlers["snow"] = Snow
}

func Snow(m *Machine) {
	m.Assert(m.Cave != nil, "Cannot snow a non-initialized field!")

	// TODO: allow a "negative" or "off" modifier to the snow command.

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
			arg1 == "(off)",
			"The only value allowed for the optional second argument is the string '(off)'. The string '%s' was provided",
			arg1,
		)
		snow.Cell = field.OffCell
	}

	snow.ApplyToField(m.Cave)
}
