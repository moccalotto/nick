package machine

func init() {
	InstructionHandlers["exit"] = Exit
}

func Exit(m *Machine) {
	m.State.PC = len(m.Tape)
}
