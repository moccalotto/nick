package machine

import (
	"github.com/moccalotto/nick/field/modifiers"
)

func init() {
	InstructionHandlers["oddities"] = Oddities
}

// oddities 50	# removes areas smaller than 50 tiles
func Oddities(m *Machine) {
	o := modifiers.NewOddityRemover(m.ArgAsInt(0))

	m.Field.Apply(o)
}
