package machine

import (
	"github.com/moccaloto/nick/field/modifiers"
	"regexp"
)

var evolver *regexp.Regexp = regexp.MustCompile(`^B(\d*)/S(\d*)$`)

func init() {
	InstructionHandlers["evolve"] = Evolve
}

func Evolve(m *Machine) {
	m.Assert(m.Field != nil, "Cannot evolve a non-initialized field!")

	arg0 := m.ArgAsString(0)
	m.Assert(
		evolver.MatchString(arg0),
		"The format '%s' is invalid. See '%s'",
		arg0,
		"https://en.wikipedia.org/wiki/Life-like_cellular_automaton#Notation_for_rules",
	)

	m.Field.Apply(modifiers.NewAutomaton(arg0))
}
