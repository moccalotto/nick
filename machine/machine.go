package machine

import (
	"fmt"
	"github.com/golang-collections/collections/stack"
	"github.com/moccalotto/nick/field"
	"math/rand"
	"strconv"
	"time"
)

type ExceptionHandler func(m *Machine, msg interface{}, a ...interface{})

type InstructionHandler func(m *Machine)

type Machine struct {
	Rng       *rand.Rand        // Random Number Generator
	Seed      int64             // Seed for the Rng
	Field     *field.Field      // field to populate.
	Stack     *stack.Stack      // Stack used for nesting and looping.
	State     *MachineState     // The current state of the machine.
	Tape      []Instruction     // the entire program.
	Trace     []int             // trace of executed instructions.
	Exception ExceptionHandler  // Exception Handler.
	Vars      map[string]string // Map of variables set inside the program.
	Limits    Restrictions      // Restrictions on runtime, cell count, etc.
	StartedAt time.Time         // When did the execution start. If nil, it hasn't started yet.
}

var InstructionHandlers map[string]InstructionHandler = make(map[string]InstructionHandler)

func (m *Machine) Assert(condition bool, msg interface{}, a ...interface{}) {
	if condition {
		return
	}

	m.Throw(msg, a...)
}

func (m *Machine) Throw(msg interface{}, a ...interface{}) {
	m.Exception(m, fmt.Sprintf("Error on line %d: %s", m.CurrentInstruction().Line, msg), a...)
}

func (m *Machine) StrToInt(s string) int {
	i, e := strconv.Atoi(s)

	m.Assert(e == nil, "Could not convert string '%s' to integer", s)

	return i
}

func (m *Machine) StrToFloat(s string) float64 {
	f, e := strconv.ParseFloat(s, 64)

	m.Assert(e == nil, "Could not convert string '%s' to float", s)

	return f
}

func (m *Machine) MustGetVar(id string) string {
	val, ok := m.Vars[id]

	m.Assert(ok, "The variable '%s' was referenced, but does not exist in map: %v", id, m.Vars)

	return val
}

func (m *Machine) MustGetString(a Arg) string {
	switch a.T {
	case StrArg, FloatArg, IntArg:
		return a.StrVal
	case CmdArg:
		// TODO: @rand, @foo, etc. could be special value handlers, just like InstructionHandlers
		// They could then be registered on runtime.
		switch a.StrVal {
		case "rand":
			return strconv.FormatFloat(m.Rng.Float64(), 'E', -1, 64)
		case "pc":
			return strconv.Itoa(m.State.PC)
		case "loop":
			return strconv.Itoa(m.State.Loop)
		default:
			m.Throw("Unknown command special command @%s (%v)", a.StrVal, a)
		}
	case VarArg:
		return m.MustGetVar(a.StrVal)
	}

	m.Throw("This should never happen MustGetString(%v)", a)

	return ""
}

func (m *Machine) MustGetFloat(a Arg) float64 {
	switch a.T {
	case StrArg:
		m.Throw("Could not convert argument %v into float", a.StrVal)
	case FloatArg:
		return a.FloatVal
	case IntArg:
		return a.FloatVal
	case CmdArg:
		// TODO: @rand, @foo, etc. could be special value handlers, just like InstructionHandlers
		// They could then be registered on runtime.
		switch a.StrVal {
		case "rand":
			return m.Rng.Float64()
		case "pc":
			return float64(m.State.PC)
		case "loop":
			return float64(m.State.Loop)
		default:
			m.Throw("Unknown command special command @%s", a.StrVal)
		}

	case VarArg:
		return m.StrToFloat(m.MustGetVar(a.StrVal))
	}

	m.Throw("This should never happen MustGetFloat(%v)", a)

	return 0.0
}

func (m *Machine) MustGetInt(a Arg) int {
	switch a.T {
	case StrArg, FloatArg:
		m.Throw("Could not convert argument (%v) into an integer", a.StrVal)
	case IntArg:
		return a.IntVal
	case CmdArg:
		// TODO: @rand, @foo, etc. could be special value handlers, just like InstructionHandlers
		// They could then be registered on runtime.
		switch a.StrVal {
		case "rand":
			m.Throw("Cannot use @rand. Expecting an integer")
		case "pc":
			return m.State.PC
		case "loop":
			return m.State.Loop
		default:
			m.Throw("Unknown command special command @%s", a.StrVal)
		}
	case VarArg:
		return m.StrToInt(m.MustGetVar(a.StrVal))
	}

	m.Throw("This should never happen MustGetInt(%v)", a)

	return 0
}

func (m *Machine) CurrentInstruction() Instruction {
	m.Assert(m.State.PC < len(m.Tape), "Program Counter is out of scope. Was the machine even loaded? (%+v)", m)

	return m.Tape[m.State.PC]
}

func (m *Machine) Arg(n int) Arg {
	instr := m.CurrentInstruction()

	m.Assert(
		n < len(instr.Args),
		"The '%s' instruction expects at least %d arguments, but %d was given!",
		instr.Cmd,
		n+1,
		len(instr.Args),
	)

	return instr.Args[n]
}

// Get then nth argument to the current instruction as a string.
// Any magic interpolations and value-replacements are done seamlessly.
func (m *Machine) ArgAsString(n int) string {
	return m.MustGetString(m.Arg(n))
}

func (m *Machine) ArgAsInt(n int) int {
	return m.MustGetInt(m.Arg(n))
}

func (m *Machine) ArgAsFloat(n int) float64 {
	return m.MustGetFloat(m.Arg(n))
}

// Number of args for the current instruction
func (m *Machine) ArgCount() int {
	return len(m.CurrentInstruction().Args)
}

// Does the current instruction have at least n+1 argmments
func (m *Machine) HasArg(n int) bool {
	return m.ArgCount() > n
}

// Push the entire state into the stack
func (m *Machine) PushState() {
	tmp := *m.State
	m.Stack.Push(tmp)
}

func (m *Machine) PopState() {
	tmp := m.Stack.Pop().(MachineState)
	m.State = &tmp
}

func (m *Machine) PeekState() *MachineState {
	tmp := m.Stack.Peek().(MachineState)
	return &tmp
}

func (m *Machine) ShouldSkip(i *Instruction) bool {
	if len(m.State.SkipUntil) == 0 {
		return false
	}

	t, ok := m.State.SkipUntil[i.Cmd]

	return t && ok
}

// Execute runs the script.
// NOTE that this can modify all of the machine's properties, except for the tape.
// If you want the machine to be pristine, you should clone the machine beforehand.
func (m *Machine) Execute() error {
	m.StartedAt = time.Now()

	for m.State.PC = 0; m.State.PC < len(m.Tape); m.State.PC++ {
		m.execCurrentInstruction()

		if err := m.checkRestrictions(); err != nil {
			return err
		}
	}

	return nil
}

func (m *Machine) execCurrentInstruction() {
	i := m.Tape[m.State.PC]

	if m.ShouldSkip(&i) {
		return
	}

	handler, ok := InstructionHandlers[i.Cmd]

	if !ok {
		m.Throw("Unknown instruction '%s' on line %d (%+v)", i.Cmd, i.Line, i)
		return
	}

	handler(m)
}
