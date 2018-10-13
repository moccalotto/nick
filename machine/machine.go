package machine

import (
	"github.com/golang-collections/collections/stack"
	"github.com/moccalotto/nick/field"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

type ExceptionHandler func(m *Machine, msg interface{}, a ...interface{})

type InstructionHandler func(m *Machine)

type Machine struct {
	Field     *field.Field      // field to populate
	Stack     *stack.Stack      // Stack used for nesting and looping
	State     *MachineState     // The current state of the machine
	Tape      []Instruction     // the entire program
	Trace     []int             // trace of executed instructions
	Exception ExceptionHandler  // Exception Handler
	Vars      map[string]string // Map of variables set inside the program
}

var InstructionHandlers map[string]InstructionHandler = make(map[string]InstructionHandler)

func (m *Machine) Assert(condition bool, msg interface{}, a ...interface{}) {
	if condition {
		return
	}

	m.Exception(m, msg, a...)
}

func (m *Machine) StrToInt(s string) int {
	i, e := strconv.Atoi(s)

	m.Assert(e == nil, "Could not convert string '%s' to ingeger: %+v", s, e)

	return i
}

func (m *Machine) StrToFloat(s string) float64 {
	f, e := strconv.ParseFloat(s, 64)

	m.Assert(e == nil, "Could not convert string '%s' to float: %+v", s, e)

	return f
}

func (m *Machine) MustGetVar(id string) string {
	val, ok := m.Vars[id]

	m.Assert(ok, "The variable '%s' was referenced, but does not exist in map: %v", id, m.Vars)

	return val
}

func (m *Machine) interpolatedValue(s string) string {
	return regexp.MustCompile(`\{\$[^$}]+\}`).ReplaceAllStringFunc(s, func(v string) string {
		var id string = v[2 : len(v)-1]

		return m.MustGetVar(id)
	})
}

func (m *Machine) MustGetString(s string) string {

	// TODO: @rand, @foo, etc. could be special value handlers, just like InstructionHandlers
	// They could then be registered on runtime.
	if s == "@rand" {
		return strconv.FormatFloat(rand.Float64(), 'E', -1, 64)
	}
	if s == "@pc" {
		return strconv.Itoa(m.State.PC)
	}
	if s == "@loop" {
		return strconv.Itoa(m.State.Loop)
	}

	if !strings.HasPrefix(s, "$") {
		return m.interpolatedValue(s)
	}

	return m.MustGetVar(s[1:])
}

func (m *Machine) MustGetFloat(s string) float64 {
	return m.StrToFloat(m.MustGetString(s))
}

func (m *Machine) MustGetInt(s string) int {
	return m.StrToInt(m.MustGetString(s))
}

func (m *Machine) CurrentInstruction() Instruction {
	m.Assert(m.State.PC < len(m.Tape), "Program Counter is out of scope. Was the machine even loaded? (%+v)", m)

	return m.Tape[m.State.PC]
}

// Get then nth argument to the current instruction as a string.
// Any magic interpolations and value-replacements are done seamlessly.
func (m *Machine) ArgAsString(n int) string {
	instr := m.CurrentInstruction()

	m.Assert(
		n < len(instr.Args),
		"The '%s' instruction expects at least %d arguments, but %d was given!",
		instr.Cmd,
		n+1,
		len(instr.Args),
	)

	return m.MustGetString(instr.Args[n])
}

func (m *Machine) ArgAsInt(n int) int {
	return m.StrToInt(m.ArgAsString(n))
}

func (m *Machine) ArgAsFloat(n int) float64 {
	s := m.ArgAsString(n)
	if strings.HasSuffix(s, "%") {
		return m.StrToFloat(s[:len(s)-1]) / 100.0
	}
	return m.StrToFloat(s)
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
	m.Stack.Push(*m.State)
	log.Printf("Pushed: %#v", m.State)
}

func (m *Machine) PopState() {
	tmp := m.Stack.Pop().(MachineState)
	m.State = &tmp
	log.Printf("Popped: %#v", m.State)
}

func (m *Machine) PeekState() *MachineState {
	tmp := m.Stack.Peek().(MachineState)
	return &tmp
}

func (m *Machine) Execute() {
	for m.State.PC = 0; m.State.PC < len(m.Tape); m.State.PC++ {
		m.execCurrentInstruction()
	}
}

func (m *Machine) execCurrentInstruction() {
	i := m.Tape[m.State.PC]

	log.Println(m.State.PC)

	handler, ok := InstructionHandlers[i.Cmd]

	if !ok {
		m.Exception(m, "Unknown instruction: %s (%+v)", i.Cmd, i)
		return
	}

	handler(m)
}
