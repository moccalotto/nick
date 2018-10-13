package machine

import "testing"

func TestInit(t *testing.T) {
	script := `init 20 x 30`

	m := MachineFromScript(script)

	Init(m)

	if m.Field.Width() != 20 {
		t.Errorf("Expected width to be 20, but it was %d", m.Field.Width())
	}
	if m.Field.Height() != 30 {
		t.Errorf("Expected height to be 20, but it was %d", m.Field.Height())
	}
}

func TestInitTooFewArgs(t *testing.T) {
	script := `init 10 20`

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

	Init(m)
}

func TestInitMalformedArg(t *testing.T) {
	script := `init 10X20`

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

	Init(m)
}

func TestInitCalledTwice(t *testing.T) {
	script := `init 10 x 20`

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

	Init(m)
	Init(m) // this time it should crash!
}
