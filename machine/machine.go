package machine

import (
	"fmt"
	"github.com/moccalotto/nick/field"
	"github.com/moccalotto/nick/utils"
	"math/rand"
	"strconv"
	"time"
)

// This is the machine, that executes cave-scripts
type Machine struct {
	Seed       int64             // Seed for the Rng.
	Rng        *rand.Rand        // Random Number Generator.
	Cave       *field.Field      // The cells/tiles of the cave.
	Stack      *Stack            // Stack used for nesting and looping.
	State      *MachineState     // The current state of the machine.
	Tape       []Instruction     // the entire program.
	Exception  ExceptionHandler  // Exception Handler.
	Vars       VarBag            // Map of variables set inside the program.
	MaxRuntime time.Duration     // Max time the machine is allowed to execute instructions.
	MaxCells   int               // Max number of cells in the cave.
	MaxWidth   int               // Max width of the cave.
	MaxHeight  int               // Max height of the cave.
	StartedAt  time.Time         // When did the execution start. If 0, it hasn't started yet.
	calculator *utils.Calculator // Math parsing engine
}

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

func (m *Machine) MustGetCmd(id string) string {
	switch id {
	case "pc":
		return strconv.Itoa(m.State.PC)
	case "loop":
		return strconv.Itoa(m.State.Loop)
	case "cond":
		if m.State.Cond {
			return "1"
		}
		return "0"
	case "line":
		return strconv.Itoa(m.CurrentInstruction().Line)
	case "width":
		return strconv.Itoa(m.Cave.Width())
	case "height":
		return strconv.Itoa(m.Cave.Height())
	default:
		defer func() {
			r := recover()
			m.Assert(
				r == nil,
				"Error '%s' during calculation [[ %s ]] on line %d.",
				r,
				id,
				m.CurrentInstruction().Line,
			)
		}()
		return strconv.FormatFloat(
			m.Calculator().Eval(id),
			'f',
			-1,
			64,
		)
		m.Throw("Unknown command special command @%s", id)
	}

	panic("Code never reached!")
}

func (m *Machine) MustGetString(a Arg) string {
	switch a.T {
	case StrArg, FloatArg, IntArg:
		return a.StrVal
	case CmdArg:
		return m.MustGetCmd(a.StrVal)
	case VarArg:
		return m.MustGetVar(a.StrVal)
	}

	m.Throw("This should never happen MustGetString(%v)", a)

	panic("Code never reached!")
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
		return m.StrToFloat(m.MustGetCmd(a.StrVal))
	case VarArg:
		return m.StrToFloat(m.MustGetVar(a.StrVal))
	}

	m.Throw("This should never happen MustGetFloat(%v)", a)

	panic("Code never reached!")
}

func (m *Machine) MustGetInt(a Arg) int {
	switch a.T {
	case StrArg, FloatArg:
		m.Throw("Could not convert argument (%v) into an integer", a.StrVal)
	case IntArg:
		return a.IntVal
	case CmdArg:
		return m.StrToInt(m.MustGetCmd(a.StrVal))
	case VarArg:
		return m.StrToInt(m.MustGetVar(a.StrVal))
	}

	panic("Code never reached!")
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
	tmp := MachineState{
		PC:        m.State.PC,
		Return:    m.State.Return,
		Loop:      m.State.Loop,
		Cond:      m.State.Cond,
		SkipUntil: m.State.SkipUntil,
	}
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

// ShouldSkip returns true if the given instruction should be skipped.
// i.e. if we should proceed to the next instruction without executing it.
func (m *Machine) ShouldSkip(i *Instruction) bool {
	// Do we want to skip any instructions at all?
	if len(m.State.SkipUntil) == 0 {
		return false
	}

	// See if i is in the list of "accepted" commands.
	t, _ := m.State.SkipUntil[i.Cmd]

	// if, and only if, t is true, then the instruction is in the whitelist,
	// and should not be skipped.
	return t == false
}

// Execute runs the script.
// NOTE that this will modify the machine's properties, except for the tape.
// If you want the machine to be pristine, you should clone the machine beforehand.
func (m *Machine) Execute() error {
	m.StartedAt = time.Now()
	m.State.Cond = true

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
