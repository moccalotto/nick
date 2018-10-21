package machine

import (
	"github.com/moccalotto/nick/effects"
)

func init() {
	InstructionHandlers["prune-small-areas"] = PruneSmallAreas
}

//  prune-small-areas 50	# removes areas smaller than 50 tiles
func PruneSmallAreas(m *Machine) {
	o := effects.NewSmallAreaCRemover(m.ArgAsInt(0))

	m.Field.Apply(o)
}
