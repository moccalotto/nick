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

type Machine struct {
	Field *field.Field  // field to populate
	Stack *stack.Stack  // Stack used for nesting and looping
	State *MachineState // The current state of the machine
	Tape  []Instruction // the entire program
	Trace []int         // trace of executed instructions
	Exception ExceptionHandler // Exception Handler
}

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
	if s == "@rand" {
		return strconv.FormatFloat(rand.Float64(), 'E', -1, 64)
	}

	if !strings.HasPrefix(s, "$") {
		return s
	}

	return m.MustGetVar(s[1:])
}
