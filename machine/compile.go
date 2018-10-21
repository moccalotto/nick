package machine

import (
	"errors"
	"fmt"
	"log"
	"math/bits"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func DefaultExceptionHandler(m *Machine, msg interface{}, a ...interface{}) {
	log.Fatalf("Exception: "+msg.(string), a...)
}

func MachineFromScript(p string) *Machine {
	seed := int64(bits.ReverseBytes64(uint64(time.Now().UnixNano())))
	source := rand.NewSource(seed)
	rng := rand.New(source)
	return &Machine{
		Rng:       rng,
		Seed:      seed,
		Stack:     NewStack(),
		State:     &MachineState{},
		Tape:      scriptToInstructions(p),
		Limits:    Restrictions{},
		Exception: DefaultExceptionHandler,
		Vars:      make(VarBag),
	}
}

func escapeString(s string) (string, []string) {
	placeholders := []string{} // placeholders

	// generate a temporary string where whitespace inside [[ ]] blocks is removed
	escaped := stringEscaper.ReplaceAllStringFunc(s, func(m string) string {
		i := len(placeholders)
		txt := m[2 : len(m)-2]

		placeholders = append(placeholders, txt)

		return fmt.Sprintf("___REP___%09d___PER___", i)
	})

	return escaped, placeholders
}

func unescapeString(s string, placeholders []string) string {
	return stringUnescaper.ReplaceAllStringFunc(s, func(m string) string {
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

// lineToInstruction generates an Instruction
// l is the line and i is the line number
func lineToInstruction(l string, i int) (*Instruction, error) {
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
		Args:    makeArgs(words[1:]),
		Comment: comment,
		Line:    i + 1,
	}, nil
}

func makeArg(w string) Arg {
	if strings.HasPrefix(w, "$") {
		return Arg{
			T:      VarArg, // this is a var-reference
			StrVal: w[1:],  // remove the $ to get the name of the var
		}
	}

	if strings.HasPrefix(w, "@") {
		return Arg{
			T:      CmdArg, // This is a "command" (i.e. a special system variable)
			StrVal: w[1:],  // remove the @ to get the name of the var
		}
	}

	// Can we parse the word into a number that fits into an int?
	// The number can be parsed a base-10 number (1234), a base-8 number (02322),
	// or a hex-number (0x4d2).
	if v, e := strconv.ParseInt(w, 0, 0); e == nil {
		return Arg{
			T:        IntArg,
			FloatVal: float64(v),
			IntVal:   int(v),
			StrVal:   w,
		}
	}

	if v, e := strconv.ParseFloat(w, 64); e == nil {
		i := int(v)
		if v == float64(i) {
			return Arg{
				T:        IntArg,
				FloatVal: v,
				IntVal:   i,
				StrVal:   strconv.FormatFloat(v, 'f', -1, 64),
			}
		}
		return Arg{
			T:        FloatArg,
			FloatVal: v,
			StrVal:   strconv.FormatFloat(v, 'f', -1, 64),
		}
	}

	if strings.HasSuffix(w, "%") {
		if v, e := strconv.ParseFloat(w[0:len(w)-1], 64); e == nil {
			v = v / 100.0
			return Arg{
				T:        FloatArg,
				FloatVal: v,
				StrVal:   strconv.FormatFloat(v, 'f', -1, 64),
			}
		}
	}

	return Arg{
		T:      StrArg,
		StrVal: w,
	}
}

func makeArgs(words []string) []Arg {
	res := make([]Arg, len(words))

	for i, w := range words {
		res[i] = makeArg(w)
	}

	return res
}

func scriptToInstructions(p string) []Instruction {
	res := []Instruction{}

	for i, s := range scriptToLines(p) {
		instr, _ := lineToInstruction(s, i)
		if instr == nil {
			continue
		}

		res = append(res, *instr)
	}

	return res
}

func scriptToLines(s string) []string {
	return lineExploder.Split(s, 65535)
}
