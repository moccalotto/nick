package machine

import (
	S "github.com/golang-collections/collections/stack"
	F "github.com/moccaloto/nick/field"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type Instruction struct {
	cmd  string
	args []string
	comment string
}

type Machine struct {
	rng   *rand.Rand    // random number generator
	field *F.Field      // field to populate
	stack *S.Stack      // stack
	pc    int           // program counter
	last  int           // last address (differs from pc-1 when inside IF statements and LOOPs)
	cond  bool          // condition bit (did last comparison succeed)
	tape  []Instruction // the entire program
	trace []Instruction // trace of executed instructions
}

func MachineFromScript(p string) *Machine {
	return &Machine{
		rng:   rand.New(rand.NewSource(time.Now().UTC().UnixNano())),
		stack: S.New(),
		pc:    0,
		last:  0,
		cond:  false,
		tape:  ScriptToInstructions(p),
	}
}

func ScriptToLines(p string) []string {
	return regexp.MustCompile(`\n+`).Split(p, 65535)
}

func LineToInstruction(l string) *Instruction {
	things := strings.SplitN(l, "#", 2)
	if len(things) == 0 {
		panic("This should never happen: Empty instruction")
	}

	words := strings.Fields(things[0])

	comment := ""
	if len(things) == 2 {
		comment = strings.TrimSpace(things[1])
	}

	if len(words) < 1 {
		panic("This should never happen: Empty words slice")
	}

	return &Instruction{
		cmd: words[0],
		args: words[1:],
		comment: comment,
	}
}

func ScriptToInstructions(p string) []Instruction {
	res := []Instruction{}

	for _, line := range ScriptToLines(p) {
		line := strings.TrimSpace(line)

		if line == "" {
			continue
		}

		res = append(res, *LineToInstruction(line))
	}

	return res
}
