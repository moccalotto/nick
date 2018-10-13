package machine

import (
	"strings"
)

func init() {
	InstructionHandlers["suggest"] = Suggest
}

// Suggest handles the "suggest" instruction.
// It adds an entry into the "suggestions" part of the machine variables.
// This essentially tells post-processors what to do, but they
// do not have to adhere to it.
func Suggest(m *Machine) {
	m.Assert(
		m.ArgAsString(1) == "=",
		"Usage: 'suggest [string] = [string]'. You did this: %s %s",
		m.CurrentInstruction().Cmd,
		strings.Join(m.CurrentInstruction().Args, " "),
	)

	varName := "suggestion." + m.ArgAsString(0)

	m.Vars[varName] = m.ArgAsString(2)
}
