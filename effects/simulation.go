package effects

import (
	"github.com/moccalotto/nick/field"
	"math/rand"
)

type Transformation struct {
	SourceFilter    map[field.Cell]bool // The transformation only runs if the source cell is represented in this map.
	Coverage        float64             // The transformation only runs if a die roll [0, 1) is lower than this number.
	TargetCell      field.Cell          // The transformation returns this cell if all matches occur.
	CheckNeighbours bool                // If true, this filter must also match the number of neighbours of the given types.
	NeighbourTypes  []field.Cell        // The transformation only runs if the number of neighbours of the given type
	NeighbourCounts [9]bool             //   is represented in NeighbourCounts.
}

type Simulation struct {
	rng             *rand.Rand
	Transformations []Transformation
}

func NewSimulation(rng *rand.Rand) *Simulation {
	return &Simulation{
		rng: rng,
	}
}

func (this *Simulation) ApplyToField(f *field.Field) {
	// Cannot run async map because of use of rng.
	f.Map(func(f *field.Field, x, y int, c field.Cell) field.Cell {

		// Lazy load cache, containing the count of neighbours to [x,y] by type.
		var countCache map[field.Cell]int

		for _, trans := range this.Transformations {
			// The transformation does not apply - die roll too low
			if trans.Coverage < 1.0 && this.rng.Float64() > trans.Coverage {
				continue
			}

			// The transformation does not match the source cell
			if !trans.SourceFilter[c] {
				continue
			}

			// The transformation does not care about the neighbours, so it passes all filters.
			if len(trans.NeighbourTypes) == 0 {
				// return now, because only one transformation can apply to a cell
				return trans.TargetCell
			}

			// The transformation wants to know about our neighbours, so we load that info into the cache.
			if len(countCache) == 0 {
				countCache = f.NeighbourCountByType(x, y)
			}

			// neighbourCount is the number of neighbours that is of one of the types
			// given in NeighbourTypes
			neighbourCount := 0

			// add up the neighbours of the matching types.
			for _, cellType := range trans.NeighbourTypes {
				neighbourCount += countCache[cellType]
			}

			// The transformation matches the number of neighbours of the given types
			if trans.NeighbourCounts[neighbourCount] {
				return trans.TargetCell
			}
		}

		return c
	})
}
