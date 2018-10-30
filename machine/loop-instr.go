package machine

var skipUntilLoop map[string]bool = map[string]bool{
	"endloop":    true,
	"loop":       true,
	"break-loop": true,
}

func init() {
	InstructionHandlers["loop"] = Loop
	InstructionHandlers["endloop"] = EndLoop
	InstructionHandlers["continue"] = EndLoop
	InstructionHandlers["break-loop"] = BreakLoop
}

func Loop(m *Machine) {
	repititions := m.ArgAsInt(0)

	m.PushState()
	if repititions == 0 {
		m.State.SkipUntil = map[string]bool{
			"endloop": true,
		}

		return
	}
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

func Continue(m *Machine) {
	m.State.Loop--
	if m.State.Loop > 0 {
		m.State.PC = m.State.Return
		return
	}

	BreakLoop(m)
}

func BreakLoop(m *Machine) {
	level := 1
	next := 0
	for i := m.State.PC; i < len(m.Tape); i++ {
		if m.Tape[i].Cmd == "loop" {
			level++
		}

		if m.Tape[i].Cmd == "endloop" {
			level--
		}

		if level == 0 {
			next = m.State.PC + i
		}
	}

	m.PopState()
	m.State.PC = next
}
