package machine

import (
	"strconv"
	"strings"
)

func init() {
	InstructionHandlers["gridwh"] = GridWH
	InstructionHandlers["grid"] = Grid
	InstructionHandlers["grid-color"] = GridColor
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

	m.Vars[".grid.cols"] = cols
	m.Vars[".grid.rows"] = rows
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

	m.Vars[".grid.width"] = width
	m.Vars[".grid.height"] = height
}

// pattern:
// grid-color rgba(255,255,255,30)
// grid-color [[#1e1e1e]]
// grid-color rgb(30,30,30)
func GridColor(m *Machine) {
	errorMsg := "Invalid arguments. Usage: 'grid-color [css color definition]'"

	m.Assert(m.HasArg(0), errorMsg)

	color := strings.ToLower(m.ArgAsString(0))

	if color == "none" {
		color = "rgba(0, 0, 0, 0)"
	}

	m.Vars[".grid.color"] = color
}
