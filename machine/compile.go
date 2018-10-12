package machine

import (
	"errors"
	"github.com/golang-collections/collections/stack"
	"github.com/op/go-logging"
	"regexp"
	"strings"
	"os"
)

var log = logging.MustGetLogger("logger")

func DefaultExceptionHandler(m *Machine, msg interface{}, a ...interface{}) {
	log.Errorf("Exception: "+msg.(string), a...)
	os.Exit(1)
}

func MachineFromScript(p string) *Machine {
	return &Machine{
		Stack: stack.New(),
		State: &MachineState{},
		Tape:  scriptToInstructions(p),
		Exception: DefaultExceptionHandler,
	}
}

func lineToInstruction(l string) (*Instruction, error) {
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

func scriptToInstructions(p string) []Instruction {
	res := []Instruction{}

	for _, line := range scriptToLines(p) {
		instr, _ := lineToInstruction(line)
		if instr == nil {
			continue
		}

		res = append(res, *instr)
	}

	return res
}

func scriptToLines(p string) []string {
	return regexp.MustCompile(`[\n\r]+`).Split(p, 65535)
}
