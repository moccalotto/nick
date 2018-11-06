package machine

const (
	// The enum values of the ArgType
	StrArg      ArgType = iota + 1 // Arg is just a string
	IntArg                         // Arg is an integer (can of course also be a float)
	FloatArg                       // Arg is a float
	VarArg                         // Arg is a var-reference (i.e. starts with $)
	CmdArg                         // Arg is a special command (i.e. starts with @)
	Calculation                    // Arg is a special command (i.e. starts with @)
)

// ArgType is an enum that is one of the StrArg, IntArg, FloatArg, CmdArg values
type ArgType int

// Arg is an argument for an instruction
type Arg struct {
	T        ArgType // The inferred type of argument
	StrVal   string  // The argument as a string (or token if you will).
	FloatVal float64 // Thi argument converted to float64
	IntVal   int     // The argument converted to int
}

// Instruction as it is stored on the machine's tape.
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
	SkipUntil InstructionFilter
}

type ExceptionHandler func(m *Machine, msg interface{}, a ...interface{})

type InstructionHandler func(m *Machine)

type InstructionFilter map[string]bool

type VarBag map[string]string
