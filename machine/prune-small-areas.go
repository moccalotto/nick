package machine

import (
	"github.com/moccalotto/nick/field/modifiers"
)

func init() {
	InstructionHandlers["prune-small-areas"] = PruneSmallAreas
}

//  prune-small-areas 50	# removes areas smaller than 50 tiles
func PruneSmallAreas(m *Machine) {
	o := modifiers.NewSmallAreaCRemover(m.ArgAsInt(0))

	m.Field.Apply(o)
}
