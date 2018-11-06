package machine

func init() {
	InstructionHandlers["background-file"] = BackgroundFile
	InstructionHandlers["background-color"] = BackgroundColor
}

// patterns:
//   background-file /path/to/my/file
func BackgroundFile(m *Machine) {
	errorMsg := "Invalid arguments. Usage: 'background-file [filename]'"

	m.Assert(m.HasArg(0), errorMsg)
	m.Vars["suggestion.background.file"] = m.ArgAsString(0)
}

// patterns:
//   background #caa
//   background #fafafa
//   background rgb(255, 123, 233)
//   background rgba(255, 123, 233, 50)
func BackgroundColor(m *Machine) {
	errorMsg := "Invalid arguments. Usage: 'background-color [css color definition]'"

	m.Assert(m.HasArg(0), errorMsg)

	m.Vars["suggestion.background.color"] = m.ArgAsString(0)
}
