package machine

import (
	"github.com/moccalotto/nick/effects"
	"github.com/moccalotto/nick/field"
	"strconv"
	"strings"
)

func init() {
	InstructionHandlers["simulation"] = StartSimulation
	InstructionHandlers["endsimulation"] = EndSimulation
	InstructionHandlers["commit"] = Commit
	InstructionHandlers["apply"] = ApplySimulation
	InstructionHandlers["map"] = Map
}

var simStacks map[*Machine][]*effects.Simulation = map[*Machine][]*effects.Simulation{}

func currentSim(m *Machine) *effects.Simulation {
	// pop the simulation off of the stack and execute it
	stack, ok := simStacks[m]

	m.Assert(ok, "No simulation started")
	m.Assert(len(stack) > 0, "Not in a simulation")

	topOpStack := len(stack) - 1

	return stack[topOpStack]
}

func StartSimulation(m *Machine) {
	newSimulation := effects.NewSimulation(m.Rng)
	simStacks[m] = append(simStacks[m], newSimulation)
}

func EndSimulation(m *Machine) {
	ApplySimulation(m)
	stack := simStacks[m]
	simStacks[m] = stack[0 : len(stack)-1]
}

func Commit(m *Machine) {
	// end + start new
	EndSimulation(m)
	StartSimulation(m)
}

func ApplySimulation(m *Machine) {
	currentSim(m).ApplyToField(m.Cave)
}

// map Alive nextTo 4,5,6,7,8 * Alive => Dead    # Death from overpopulation
func Map(m *Machine) {

	// arg: type
	// 0: SourceFilter
	// 1: "nextTo"
	// 2: NeighbourCounts
	// 3: "*"
	// 4: NeighbourTypes
	// 5: "=>"
	// 6: TargetCell

	m.Assert(m.ArgAsString(1) == "nextTo", "second arg must be 'nextTo'")
	m.Assert(m.ArgAsString(3) == "*", "fourth arg must be '*'")
	m.Assert(m.ArgAsString(5) == "=>", "sixth arg must be '=>'")

	t := effects.Transformation{}

	if len(simStacks[m]) == 0 {
		m.Assert(len(simStacks[m]) > 0, "Simulation not started")
	}

	topOfStack := len(simStacks[m]) - 1

	simStacks[m][topOfStack].Transformations = append(
		simStacks[m][topOfStack].Transformations,
		t,
	)

	sourcesStr := m.ArgAsString(0)
	t.SourceFilter = map[field.Cell]bool{}
	t.Coverage = 1.0
	for _, cellAlias := range strings.Split(sourcesStr, ",") {
		if cell, ok := m.CellNames[cellAlias]; ok {
			t.SourceFilter[cell] = true
		} else if num, err := strconv.Atoi(cellAlias); err == nil {
			t.SourceFilter[field.Cell(num)] = true
		} else {
			panic(err)
		}
	}

	nCountStr := m.ArgAsString(2)

	for _, numStr := range strings.Split(nCountStr, ",") {
		if i, err := strconv.Atoi(numStr); err == nil {
			t.NeighbourCounts[i] = true
		} else {
			panic(err)
		}
	}

	neighboursStr := m.ArgAsString(4)
	t.NeighbourTypes = []field.Cell{}
	t.Coverage = 1.0
	for _, cellAlias := range strings.Split(neighboursStr, ",") {
		if cell, ok := m.CellNames[cellAlias]; ok {
			t.NeighbourTypes = append(t.NeighbourTypes, cell)
		} else if num, err := strconv.Atoi(cellAlias); err == nil {
			t.NeighbourTypes = append(t.NeighbourTypes, field.Cell(num))
		} else {
			panic(err)
		}
	}

	targetCellStr := m.ArgAsString(6)

	if cell, ok := m.CellNames[targetCellStr]; ok {
		t.TargetCell = cell
	} else {
		if num, err := strconv.Atoi(targetCellStr); err == nil {
			t.TargetCell = field.Cell(num)
		} else {
			panic(err)
		}
	}
}
