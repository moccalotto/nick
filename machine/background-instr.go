package machine

import (
	"strings"
)

func init() {
	InstructionHandlers["background-file"] = BackgroundFile
	InstructionHandlers["background-color"] = BackgroundColor
}

// patterns:
//   background-file /path/to/my/file
func BackgroundFile(m *Machine) {
	errorMsg := "Invalid arguments. Usage: 'background-file [filename]'"

	m.Assert(m.HasArg(0), errorMsg)
	m.Vars[".background.file"] = m.ArgAsString(0)
}

// patterns:
//   background-color #caa
//   background-color #fafafa
//   background-color rgb(255, 123, 233)
//   background-color rgba(255, 123, 233, 50)
func BackgroundColor(m *Machine) {
	errorMsg := "Invalid arguments. Usage: 'background-color [css color definition]'"

	m.Assert(m.HasArg(0), errorMsg)

	color := strings.ToLower(m.ArgAsString(0))

	if color == "none" {
		color = "rgba(0, 0, 0, 0)"
	}

	m.Vars[".background.color"] = color
}
