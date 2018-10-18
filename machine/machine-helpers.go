package machine

import (
	"fmt"
	"time"
)

const (
	// The enum values of the ArgType
	StrArg   ArgType = 1 // Arg is just a string
	IntArg   ArgType = 2 // Arg is an integer (can of course also be a float)
	FloatArg ArgType = 3 // Arg is a float
	VarArg   ArgType = 4 // Arg is a var-reference (i.e. starts with $)
	CmdArg   ArgType = 5 // Arg is a special command (i.e. starts with @)
)

// ArgType is an enum that is one of the StrArg, IntArg, FloatArg, CmdArg values
type ArgType int

// Arg is an argument for an instruction
type Arg struct {
	T        ArgType // The inferred type of argument
	StrVal   string  // The argument as a string
	FloatVal float64 // Thi argument converted to float64
	IntVal   int     // The argument converted to int
}

// Instruction is an instruction for a Machine to perform.
type Instruction struct {
	Cmd     string // The name of this instruction (for instance "init")
	Args    []Arg  // A list of arguments for the instruction
	Comment string // The comment associated with this instruction (if any was given)
	Line    int    // The line number in the script
}

// This entire state is pushed whenever we enter a control structure
type MachineState struct {
	PC        int  // program counter
	Return    int  // Return address (to return to in loops, if-branches, etc.)
	Loop      int  // Loop Counter (used to count iterations inside iterators)
	Cond      bool // condition bit (did last comparison succeed)
	SkipUntil map[string]bool
}

type Restrictions struct {
	MaxRuntime time.Duration
	MaxCells   int
	MaxWidth   int
	MaxHeight  int
}

// checkRestrictions returns an error if m.Limits is not adhered to
func (m *Machine) checkRestrictions() error {
	if err := m.timedOut(); err != nil {
		return err
	}

	if err := m.tooManyCells(); err != nil {
		return err
	}

	if err := m.tooWide(); err != nil {
		return err
	}

	if err := m.tooTall(); err != nil {
		return err
	}

	return nil
}

// timeOut returns an error if a machine's execution has taken too long.
// the time Limits is given by Machine.MaxRunTime
func (m *Machine) timedOut() error {
	// no max time specified, i.e. we can go on forever.
	if m.Limits.MaxRuntime == 0 {
		return nil
	}

	runtime := time.Now().Sub(m.StartedAt)

	if runtime > m.Limits.MaxRuntime {
		return fmt.Errorf("Timed out after %f seconds", runtime.Seconds())
	}

	return nil
}

// tooManyCells returns an error if the field has too many cells
// the Limits on cell count is given in Machine.MaxCells
func (m *Machine) tooManyCells() error {
	// we don't have a field yet, so it can't be too large
	if m.Field == nil {
		return nil
	}

	w := m.Field.Width()
	h := m.Field.Height()
	max := m.Limits.MaxCells

	// there is no maximum number of cells
	if max <= 0 {
		return nil
	}

	if w*h > max {
		return fmt.Errorf(
			"Grid is too large. Max number of cells allowed is %d, but the current size is (%dx%d) %d",
			max,
			w,
			h,
			w*h,
		)
	}

	return nil
}

func (m *Machine) tooWide() error {
	// we don't have a field yet, so it can't be too large
	if m.Field == nil {
		return nil
	}

	w := m.Field.Width()
	max := m.Limits.MaxWidth

	// there is no maximum width
	if max <= 0 {
		return nil
	}

	if w > max {
		return fmt.Errorf(
			"Grid is too large. Max width is %d, but the current width is %d",
			max,
			w,
		)
	}

	return nil
}

func (m *Machine) tooTall() error {
	// we don't have a field yet, so it can't be too large
	if m.Field == nil {
		return nil
	}

	h := m.Field.Height()
	max := m.Limits.MaxHeight

	// there is no maximum height
	if max <= 0 {
		return nil
	}

	if h > max {
		return fmt.Errorf(
			"Grid is too tall. Max height is %d, but the current height is %d",
			max,
			h,
		)
	}

	return nil
}
