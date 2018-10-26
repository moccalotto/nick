package machine

import "github.com/moccalotto/nick/effects"

func init() {
	InstructionHandlers["gridn"] = GridN
	InstructionHandlers["grid"] = Grid
}

// pattern:
// gridn [cols] x [rows]
// OR
// gridn [n]
func GridN(m *Machine) {
	cols := m.ArgAsInt(0)
	rows := cols
	errorMsg := "Invalid arguments. Usage: 'grid [int]' or 'grid [cols] x [rows]'"
	if m.ArgCount() == 3 {
		m.Assert(m.ArgAsString(1) == "x", errorMsg)
		rows = m.ArgAsInt(2)
	}

	grid := effects.NewGridNM(cols, rows)

	grid.ApplyToField(m.Field)
}

// pattern:
// gridn [width] x [height]
// OR
// gridn [n]
func Grid(m *Machine) {
	width := m.ArgAsInt(0)
	height := width
	errorMsg := "Invalid arguments. Usage: 'grid [int]' or 'grid [width] x [height]'"
	if m.ArgCount() == 3 {
		m.Assert(m.ArgAsString(1) == "x", errorMsg)
		height = m.ArgAsInt(2)
	}
	grid := effects.NewGrid(width, height)

	grid.ApplyToField(m.Field)
}
