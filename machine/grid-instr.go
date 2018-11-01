package machine

import (
	"strconv"
)

func init() {
	InstructionHandlers["gridwh"] = GridWH
	InstructionHandlers["grid"] = Grid
}

// pattern:
// grid [cols] x [rows]
// OR
// grid [n]
func Grid(m *Machine) {
	errorMsg := "Invalid arguments. Usage: 'grid [int]' or 'grid [cols] x [rows]'"

	cols := m.ArgAsString(0)
	rows := cols
	if m.ArgCount() == 3 {
		m.Assert(m.ArgAsString(1) == "x", errorMsg)
		rows = m.ArgAsString(2)
	}

	if _, err := strconv.Atoi(cols); err != nil {
		m.Throw(errorMsg)
	}
	if _, err := strconv.Atoi(rows); err != nil {
		m.Throw(errorMsg)
	}

	m.Vars["suggestion.grid.cols"] = cols
	m.Vars["suggestion.grid.rows"] = rows
}

// pattern:
// gridWH [width] x [height]
// OR
// gridWH [n]
func GridWH(m *Machine) {
	errorMsg := "Invalid arguments. Usage: 'gridwh [int]' or 'gridwh [width] x [height]'"

	width := m.ArgAsString(0)
	height := width

	if m.ArgCount() == 3 {
		m.Assert(m.ArgAsString(1) == "x", errorMsg)
		height = m.ArgAsString(2)
	}

	if _, err := strconv.Atoi(width); err != nil {
		m.Throw(errorMsg)
	}
	if _, err := strconv.Atoi(height); err != nil {
		m.Throw(errorMsg)
	}

	m.Vars["suggestion.grid.width"] = width
	m.Vars["suggestion.grid.height"] = height
}
