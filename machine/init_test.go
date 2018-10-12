package machine

import "testing"

func TestHandleInitInstructionOK(t *testing.T) {
	script = `init 20x30`

	m := MachineFromScript(script)

	HandleInitInstruction(m, m.Tape[0].Args)

	if m.Field.Width() != 20 {
		t.Errorf("Expected width to be 20, but it was %d", m.Field.Width())
	}
	if m.Field.Height() != 30 {
		t.Errorf("Expected height to be 20, but it was %d", m.Field.Height())
	}
}

func TestHandleInitInstructionTooManyArgs(t *testing.T) {
	script = `init 10 20`

	m := MachineFromScript(script)

	m.Exception = func(m *Machine, msg interface{}, a ...interface{}) {
		panic(msg.(string))
	}
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Expected the script '%s' to throw an exception, but it didn't!", script)
		} else {
			t.Log("We got this error fancy error: " + r.(string))
		}
	}()

	HandleInitInstruction(m, m.Tape[0].Args)
}

func TestHandleInitInstructionMalformedArg(t *testing.T) {
	script = `init 10X20`

	m := MachineFromScript(script)

	m.Exception = func(m *Machine, msg interface{}, a ...interface{}) {
		panic(msg.(string))
	}
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Expected the script '%s' to throw an exception, but it didn't!", script)
		} else {
			t.Log("We got this error fancy error: " + r.(string))
		}
	}()

	HandleInitInstruction(m, m.Tape[0].Args)
}

func TestHandleInitInstructionCalledTwice(t *testing.T) {
	script = `init 10x20`

	m := MachineFromScript(script)

	m.Exception = func(m *Machine, msg interface{}, a ...interface{}) {
		panic(msg.(string))
	}
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Expected the script '%s' to throw an exception, but it didn't!", script)
		} else {
			t.Log("We got this error fancy error: " + r.(string))
		}
	}()

	HandleInitInstruction(m, m.Tape[0].Args)
	HandleInitInstruction(m, m.Tape[0].Args) // this time it should crash!
}
