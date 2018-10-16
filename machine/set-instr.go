package machine

import (
	"strconv"
)

func init() {
	InstructionHandlers["set"] = Set
	InstructionHandlers["set-rand-int"] = SetRandInt
	InstructionHandlers["set-rand-float"] = SetRandFloat
}

// set $foo = [[this is a long string]]
func Set(m *Machine) {
	a := m.Arg(0)
	errStr := "Invalid use of instruction. Correct use is: set $varname = [value]"
	m.Assert(a.T == VarArg, errStr)

	m.Assert(m.ArgAsString(1) == "=", "Invalid use of instruction. Correct use is set $varname = [value]")

	m.Vars[a.StrVal] = m.ArgAsString(2)
}

// set-rand-int $width  = 1 to 3
// pattern: set-rand-int {varname:var} = {from:int} to {to:int}
func SetRandInt(m *Machine) {
	errStr := "Invalid use of instruction. Correct use is: set-rand-int $varname = [min] to [max]"
	a := m.Arg(0)
	m.Assert(a.T == VarArg, errStr)

	m.Assert(m.ArgAsString(1) == "=", errStr)
	m.Assert(m.ArgAsString(3) == "to", errStr)

	min := m.ArgAsInt(2)
	max := m.ArgAsInt(4)
	delta := max - min

	m.Vars[a.StrVal] = strconv.Itoa(m.Rng.Intn(delta) + min)
}

// set-rand-int $width  = 1 to 3
// pattern: set-rand-int {varname:var} = {from:int} to {to:int}
func SetRandFloat(m *Machine) {
	errStr := "Invalid use of instruction. Correct use is: set-rand-float $varname = [min] to [max]"
	a := m.Arg(0)
	m.Assert(a.T == VarArg, errStr)

	m.Assert(m.ArgAsString(1) == "=", errStr)
	m.Assert(m.ArgAsString(3) == "to", errStr)

	min := m.ArgAsFloat(2)
	max := m.ArgAsFloat(4)
	delta := max - min

	f := m.Rng.Float64()*delta + min

	m.Vars[a.StrVal] = strconv.FormatFloat(f, 'f', -1, 64)
}
