package machine

import (
	"errors"
	"fmt"
	"github.com/golang-collections/collections/stack"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func DefaultExceptionHandler(m *Machine, msg interface{}, a ...interface{}) {
	log.Fatalf("Exception: "+msg.(string), a...)
}

func MachineFromScript(p string) *Machine {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	rng := rand.New(source)
	return &Machine{
		Rng:       rng,
		Seed:      seed,
		Stack:     stack.New(),
		State:     &MachineState{},
		Tape:      scriptToInstructions(p),
		Limits:    Restrictions{},
		Exception: DefaultExceptionHandler,
		Vars:      map[string]string{},
	}
}

func escapeString(s string) (string, []string) {
	r := regexp.MustCompile(`\[\[[^\]]+\]\]`)
	placeholders := []string{} // placeholders

	// generate a temporary string where whitespace inside [[ ]] blocks is removed
	escaped := r.ReplaceAllStringFunc(s, func(m string) string {
		i := len(placeholders)
		txt := m[2 : len(m)-2]

		placeholders = append(placeholders, txt)

		return fmt.Sprintf("___REP___%09d___PER___", i)
	})

	return escaped, placeholders
}

func unescapeString(s string, placeholders []string) string {
	r := regexp.MustCompile(`___REP___\d{9}___PER___`)

	return r.ReplaceAllStringFunc(s, func(m string) string {
		var i int
		_, err := fmt.Sscanf(m, "___REP___%09d___PER___", &i)

		if err != nil {
			return m
		}

		return placeholders[i]
	})
}

func extractWords(s string, placeholders []string) []string {
	// split the temporary string into words
	splitted := strings.Fields(s)

	res := make([]string, len(splitted))

	// Replace the placeholders in each word with their original meaning
	for i, word := range splitted {
		res[i] = unescapeString(word, placeholders)
	}

	return res
}

func extractCodeAndComment(s string) (string, string) {
	parts := strings.SplitN(s, "#", 2)

	code := strings.TrimSpace(strings.Replace(parts[0], "\\#", "#", -1))
	comment := ""
	if len(parts) == 2 {
		comment = strings.TrimSpace(parts[1])
	}

	return code, comment
}

func lineToInstruction(l string) (*Instruction, error) {
	escaped, placeholders := escapeString(l)
	code, comment := extractCodeAndComment(escaped)
	if len(code) == 0 {
		// string was empty or contained only whitespace
		return nil, errors.New("Empty Code Block")
	}

	words := extractWords(code, placeholders)
	comment = unescapeString(comment, placeholders)

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
