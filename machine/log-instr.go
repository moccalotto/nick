package machine

import "strings"
import L "log"

func init() {
	InstructionHandlers["log"] = Log
}

// TODO: this should log via the machine. The machine should have a mechanism for outputting log commands.
func Log(m *Machine) {
	instr := m.CurrentInstruction()
	L.Printf("LOG: %s (%+v)", strings.Join(instr.Args, " "), m)
}
