package machine

func init() {
	InstructionHandlers["stop!"] = Stop
}

func Stop(m *Machine) {
	m.State.PC = len(m.Tape)
}
