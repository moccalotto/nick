package machine

import (
	"github.com/golang-collections/collections/stack"
	"github.com/moccaloto/nick/field"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

type ExceptionHandler func(m *Machine, msg interface{}, a ...interface{})

type InstructionHandler func(m *Machine)

type Machine struct {
	Field     *field.Field     // field to populate
	Stack     *stack.Stack     // Stack used for nesting and looping
	State     *MachineState    // The current state of the machine
	Tape      []Instruction    // the entire program
	Trace     []int            // trace of executed instructions
	Exception ExceptionHandler // Exception Handler
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
	if val, ok := m.State.Vars[id]; ok {
		return val
	}
	panic("DAFUQ MAYN: var " + id + " does not exist!")
}

func (m *Machine) getInterpolatedValue(s string) string {
	return regexp.MustCompile(`\{\$[^$)+\}`).ReplaceAllStringFunc(s, func(v string) string {
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

	if !strings.HasPrefix(s, "$") {
		return s
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
	// TODO:
	// Consider allowing percentages.
	// For instance: 85% would be converted to 0.85
	return m.StrToFloat(m.ArgAsString(n))
}

func (m *Machine) Execute() {
	for m.State.PC = 0; m.State.PC < len(m.Tape); m.State.PC++ {
		m.execCurrentInstruction()
	}
}

// Does the current instruction have at least n+1 argmments
func (m *Machine) HasArg(n int) bool {
	instr := m.CurrentInstruction()

	return len(instr.Args) > n
}

func (m *Machine) execCurrentInstruction() {
	i := m.Tape[m.State.PC]

	handler, ok := InstructionHandlers[i.Cmd]

	if !ok {
		m.Exception(m, "Unknown instruction: %s (%+v)", i.Cmd, i)
		return
	}

	handler(m)
}
