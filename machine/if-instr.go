package machine

func init() {
	InstructionHandlers["if"] = If
	InstructionHandlers["elsif"] = Elsif
	InstructionHandlers["else"] = Else
	InstructionHandlers["endif"] = Endif
}

func getCondition(m *Machine) bool {
	switch operator := m.ArgAsString(1); operator {
	case "==":
		return m.ArgAsString(0) == m.ArgAsString(2)
	case ">":
		return m.ArgAsFloat(0) > m.ArgAsFloat(2)
	case "<":
		return m.ArgAsFloat(0) < m.ArgAsFloat(2)
	case ">=":
		return m.ArgAsFloat(0) >= m.ArgAsFloat(2)
	case "<=":
		return m.ArgAsFloat(0) <= m.ArgAsFloat(2)
	case "!=":
		return m.ArgAsFloat(0) != m.ArgAsFloat(2)
	default:
		m.Throw("Unknown operator '%s', operator")
	}

	return false
}

// pattern:
// if [value] [operator] [value]
func If(m *Machine) {
	// we always push the current state when we reach an if-block
	m.PushState()

	// are we already inside a "failed" if-block?
	// if so, we don't do anything.
	// the previously pushed state will be
	// popped when we reach the endif later.
	if m.State.Cond == false {
		return
	}

	m.State.Cond = getCondition(m)

	if m.State.Cond == true {
		// Allow the if-body to be executed.
		m.State.SkipUntil = map[string]bool{}
		return
	}

	// the condition failed, we skip until we reach endif, else or elsif
	m.State.SkipUntil = map[string]bool{
		"if":    true,
		"else":  true,
		"endif": true,
		"elsif": true,
	}

	return
}

func Else(m *Machine) {

	if m.State.Cond {
		// we've just successfully executed an IF instruction.
		// therefore we do not execute the else-instruction
		m.State.Cond = false

		return
	}

	// revert condition bit to where we were before the IF statement.
	// and allow the else-body to be executed
	m.State.Cond = m.PeekState().Cond
}

func Elsif(m *Machine) {
	if m.State.Cond {
		// we've just successfully executed an IF instruction.
		// therefore we do not execute the elsif-instruction
		m.State.Cond = false

		return
	}

	m.PushState()

	m.State.Cond = getCondition(m)

	if m.State.Cond == true {
		// Allow the elsif-body to be executed.
		m.State.SkipUntil = map[string]bool{}
	}
}

func Endif(m *Machine) {
	// the if-logic is done.
	//  revert to the state before the IF statement started,
	// but make sure the Program Counter is not reverted.
	pc := m.State.PC
	m.PopState()
	m.State.PC = pc
}
