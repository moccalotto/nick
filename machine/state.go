package machine

import (
	"errors"
	S "github.com/golang-collections/collections/stack"
	F "github.com/moccaloto/nick/field"
	"regexp"
	"strings"
)

type Instruction struct {
	Cmd     string
	Args    []string
	Comment string
}

type Machine struct {
	Field *F.Field      // field to populate
	Stack *S.Stack      // Stack used for nesting and looping
	PC    int           // program counter
	Last  int           // Last address (differs from PC-1 when inside IF statements and LOOPs)
	Cond  bool          // condition bit (did last comparison succeed)
	Tape  []Instruction // the entire program
	Trace []int         // trace of executed instructions
}

func MachineFromScript(p string) *Machine {
	return &Machine{
		Stack: S.New(),
		PC:    0,
		Last:  0,
		Cond:  false,
		Tape:  ScriptToInstructions(p),
	}
}

func ScriptToLines(p string) []string {
	return regexp.MustCompile(`[\n\r]+`).Split(p, 65535)
}

func LineToInstruction(l string) (*Instruction, error) {
	l = strings.TrimSpace(l)
	parts := strings.SplitN(l, "#", 2)
	if len(parts) == 0 {
		// string was empty or contained only whitespace
		return nil, errors.New("Empty input string")
	}

	words := strings.Fields(parts[0])

	comment := ""
	if len(parts) == 2 {
		comment = strings.TrimSpace(parts[1])
	}

	if len(words) < 1 {
		// empty instruction (i.e. a line with only a comment) yields no Instruction
		return nil, errors.New("Empty instruction")
	}

	return &Instruction{
		Cmd:     words[0],
		Args:    words[1:],
		Comment: comment,
	}, nil
}

func ScriptToInstructions(p string) []Instruction {
	res := []Instruction{}

	for _, line := range ScriptToLines(p) {
		instr, _ := LineToInstruction(line)
		if instr == nil {
			continue
		}

		res = append(res, *instr)
	}

	return res
}
