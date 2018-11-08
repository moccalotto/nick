package machine

import (
	"strings"
)

func init() {
	InstructionHandlers["wall-color"] = WallColor
}

// patterns:
//   wall-color #caa
//   wall-color #fafafa
//   wall-color rgb(255, 123, 233)
//   wall-color rgba(255, 123, 233, 50)
func WallColor(m *Machine) {
	errorMsg := "Invalid arguments. Usage: 'wall-color [css color definition]'"

	m.Assert(m.HasArg(0), errorMsg)

	color := strings.ToLower(m.ArgAsString(0))

	if color == "none" {
		color = "rgba(0, 0, 0, 0)"
	}

	m.Vars[".wall.color"] = color
}
