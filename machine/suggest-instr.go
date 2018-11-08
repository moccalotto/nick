package machine

func init() {
	InstructionHandlers["suggest"] = Suggest
}

// Suggest handles the "suggest" instruction.
// It adds an entry into the "s" part of the machine variables.
// This essentially tells post-processors what to do, but they
// do not have to adhere to it.
func Suggest(m *Machine) {
	m.Assert(m.ArgAsString(1) == "=", "Usage: 'suggest [string] = [string]'")

	varName := "." + m.ArgAsString(0)

	m.Vars[varName] = m.ArgAsString(2)
}
