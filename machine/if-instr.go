package machine

func init() {
	InstructionHandlers["if"] = If
	InstructionHandlers["elsif"] = Elsif
	InstructionHandlers["else"] = Else
	InstructionHandlers["endif"] = Endif
}

func getCondition(m *Machine) bool {
	condition := false
	switch operator := m.ArgAsString(1); operator {
	case "==":
		condition = m.ArgAsString(0) == m.ArgAsString(2)
	case ">":
		condition = m.ArgAsFloat(0) > m.ArgAsFloat(2)
	case "<":
		condition = m.ArgAsFloat(0) < m.ArgAsFloat(2)
	case ">=":
		condition = m.ArgAsFloat(0) >= m.ArgAsFloat(2)
	case "<=":
		condition = m.ArgAsFloat(0) <= m.ArgAsFloat(2)
	case "!=":
		condition = m.ArgAsFloat(0) != m.ArgAsFloat(2)
	default:
		m.Throw("Unknown operator '%s', operator")
	}

	return condition
}

func If(m *Machine) {
	// pattern:
	// if [value] [operator] [value]

	m.State.Cond = getCondition(m)

	m.PushState()
	if m.State.Cond == false {
		// the condition failed, we skip until we reach endif, else or elsif
		m.State.SkipUntil = map[string]bool{
			"else":  true,
			"endif": true,
			"elsif": true,
		}

		return
	}

	// Allow the if-body to be executed.
	m.State.SkipUntil = map[string]bool{}
}

func Else(m *Machine) {

	if m.State.Cond {
		// we've just successfully executed an IF instruction.
		// therefore we do not execute the elsif (nor an else for that matter)

		m.State.Cond = false
		m.State.SkipUntil = map[string]bool{
			"endif": true,
		}

		return
	}

	// revert condition bit to where we were before the IF statement.
	// and allow the else-body to be executed
	m.State.Cond = m.PeekState().Cond
}

func Elsif(m *Machine) {
	if m.State.Cond {
		// we've just successfully executed an IF instruction.
		// therefore we do not execute the elsif (nor an else for that matter)
		m.State.Cond = false
		m.State.SkipUntil = map[string]bool{
			"endif": true,
		}
		return
	}

	m.State.Cond = getCondition(m)

	if m.State.Cond == false {
		// the condition failed, we skip until we reach endif, else or elsif
		m.State.SkipUntil = map[string]bool{
			"else":  true,
			"endif": true,
			"elsif": true,
		}

		return
	}

	// Allow the elsif-body to be executed.
	m.State.SkipUntil = map[string]bool{}
}

func Endif(m *Machine) {
	// the if-logic is done.
	//  revert to the state before the IF statement started,
	// but make sure the Program Counter is not reverted.
	pc := m.State.PC
	m.PopState()
	m.State.PC = pc
}
