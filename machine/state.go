package machine

const (
	NONE = 0
	LOOP = 1
	IF   = 2
)

type Instruction struct {
	Cmd     string
	Args    []string
	Comment string
}

// This entire state is pushed whenever we enter a control structure
type MachineState struct {
	PC     int  // program counter
	Return int  // Return address (to return to in loops, if-branches, etc.)
	Loop   int  // Loop Counter (used to count iterations inside iterators)
	Cond   bool // condition bit (did last comparison succeed)
}
