package modifiers

import (
	"github.com/moccalotto/nick/field"
	"regexp"
	"strconv"
	"strings"
)

// Automaton for evolving a field.
type Automaton struct {
	// lookup tables
	B [9]bool // if B[X] is true, it means that a cell will be born if it has X neighbours
	S [9]bool // if S[X] is true, it means that a cell will survive if it has X neighbours
}

// Create new Automaton
func NewAutomatonBool(b, s [9]bool) *Automaton {
	return &Automaton{b, s}
}

func NewAutomaton(str string) *Automaton {
	_s := [9]bool{}
	_b := [9]bool{}

	re := regexp.MustCompile("^B([0-9]*)/S([0-9]*)$")
	matches := re.FindStringSubmatch(str)
	if len(matches) != 3 {
		// TODO: correct error handling.
		panic("Bad Rule-String")
	}

	bDigits := strings.Split(matches[1], "")
	sDigits := strings.Split(matches[2], "")

	for _, digit := range bDigits {
		v, e := strconv.Atoi(digit)
		if e != nil {
			panic(e)
		}
		_b[v] = true
	}
	for _, digit := range sDigits {
		v, e := strconv.Atoi(digit)
		if e != nil {
			panic(e)
		}
		_s[v] = true
	}

	return &Automaton{_b, _s}
}

// Do the rules allow survival for the given neighbour count?
func (ca *Automaton) Survival(neighbourCount int) bool {
	return ca.S[neighbourCount]
}

// Do the rules allow giving birth for the given neighbour count?
func (ca *Automaton) Birth(neighbourCount int) bool {
	return ca.B[neighbourCount]
}

func (ca *Automaton) String() string {
	var buf strings.Builder
	buf.WriteString("B")
	for i := 0; i <= 8; i++ {
		if ca.B[i] {
			buf.WriteString(strconv.Itoa(i))
		}
	}
	buf.WriteString("/S")
	for i := 0; i <= 8; i++ {
		if ca.S[i] {
			buf.WriteString(strconv.Itoa(i))
		}
	}
	return buf.String()
}

// Next returns the state of the specified cell at the next time step.
func (ca *Automaton) NextCellState(f *field.Field, x, y int) bool {
	neighbourCount := f.NeighbourCount(x, y)
	if f.Alive(x, y) {
		return ca.Survival(neighbourCount)
	}

	return ca.Birth(neighbourCount)
}

// Apply a CA to the field
func (ca *Automaton) ApplyToField(f *field.Field) {
	tmp := field.NewField(f.Width(), f.Height())

	for x := 0; x < f.Width(); x++ {
		for y := 0; y < f.Height(); y++ {
			tmp.SetAlive(x, y, ca.NextCellState(f, x, y))
		}
	}

	f.SetCells(tmp.Cells())
}
