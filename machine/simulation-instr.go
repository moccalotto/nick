package machine

import (
	"github.com/moccalotto/nick/effects"
	"github.com/moccalotto/nick/field"
)

func init() {
	InstructionHandlers["simulation"] = Simulation
	InstructionHandlers["endsimulation"] = EndSimulation
	InstructionHandlers["commit"] = Commit
}

var simQueues map[*Machine][]Simulation

func Simulation(m *Machine) {
	// SourceFilter    map[field.Cell]bool // The transformation only runs if the source cell is represented in this map.
	// Coverage        float64             // The transformation only runs if a die roll [0, 1) is lower than this number.
	// TargetCell      field.Cell          // The transformation returns this cell if all matches occur.
	// CheckNeighbours bool                // If true, this filter must also match the number of neighbours of a particular type.
	// NeighbourType   field.Cell          // The transformation only runs if the number of neighbours of the given type
	// NeighbourCounts [9]bool             //   is represented in NeighbourCounts.

	newSimulation := effects.Simulation{
		SourceFilter: map[field.Cell]bool{},
		Coverage:     1.0,
	}
	simQueues[m] = append(newSimulation, simQueues)
}

func EndSimulation(m *Machine) {
	// pop the simulation off of the queue and execute it
}

func Commit(m *Machine) {
	// end + start new
}
