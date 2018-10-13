package machine

func init() {
	InstructionHandlers["loop"] = Loop
	InstructionHandlers["endloop"] = EndLoop
}

func Loop(m *Machine) {
	repititions := m.ArgAsInt(0)

	m.Assert(repititions > 0, "Number of repotitions must be > 0")
	m.PushState()
	m.State.Loop = repititions
	m.State.Return = m.State.PC
}

func EndLoop(m *Machine) {
	m.State.Loop--
	if m.State.Loop > 0 {
		m.State.PC = m.State.Return
		return
	}
	next := m.State.PC

	m.PopState()
	m.State.PC = next
}
