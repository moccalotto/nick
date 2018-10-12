package machine

type Instruction struct {
	Cmd     string
	Args    []string
	Comment string
}

// This entire state is pushed whenever we enter a control structure
type MachineState struct {
	PC   int               // program counter
	Last int               // Last address (differs from PC-1 when inside IF statements and LOOPs)
	Loop int               // Loop Counter (used to count iterations inside iterators)
	Cond bool              // condition bit (did last comparison succeed)
	Vars map[string]string // Map of variables set inside the program
}
