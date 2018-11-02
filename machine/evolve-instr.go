package machine

import (
	"github.com/moccalotto/nick/effects"
	"regexp"
)

var evolver *regexp.Regexp = regexp.MustCompile(`^B(\d*)/S(\d*)$`)

func init() {
	InstructionHandlers["evolve"] = Evolve
}

func Evolve(m *Machine) {
	m.Assert(m.Cave != nil, "The '%s' instruction can only be used after the 'init'", m.CurrentInstruction().Cmd)

	arg0 := m.ArgAsString(0)
	m.Assert(
		evolver.MatchString(arg0),
		"The format '%s' is invalid. See '%s'",
		arg0,
		"https://en.wikipedia.org/wiki/Life-like_cellular_automaton#Notation_for_rules",
	)

	effects.NewAutomaton(arg0).ApplyToField(m.Cave)
}
