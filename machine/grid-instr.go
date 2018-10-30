package machine

import (
	"strconv"
)

func init() {
	InstructionHandlers["gridnm"] = GridNM
	InstructionHandlers["grid"] = Grid
}

// pattern:
// gridnm [cols] x [rows]
// OR
// gridnm [n]
func GridNM(m *Machine) {
	errorMsg := "Invalid arguments. Usage: 'gridnm [int]' or 'gridnm [cols] x [rows]'"

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
// grid [width] x [height]
// OR
// grid [n]
func Grid(m *Machine) {
	errorMsg := "Invalid arguments. Usage: 'grid [int]' or 'grid [width] x [height]'"

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
