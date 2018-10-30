package effects

import (
	"github.com/moccalotto/nick/field"
	"regexp"
	"strconv"
	"strings"
)

var automatonParser *regexp.Regexp = regexp.MustCompile("^B([0-9]*)/S([0-9]*)$")

// Automaton for evolving a field.
type Automaton struct {
	// lookup tables
	B [9]field.Cell // if B[X] is true, it means that a cell will be born if it has X neighbours
	S [9]field.Cell // if S[X] is true, it means that a cell will survive if it has X neighbours
}

// Create new Automaton
func NewAutomatonBool(b, s [9]field.Cell) *Automaton {
	return &Automaton{b, s}
}

func NewAutomaton(str string) *Automaton {
	_s := [9]field.Cell{}
	_b := [9]field.Cell{}

	automatonParser := regexp.MustCompile("^B([0-9]*)/S([0-9]*)$")
	matches := automatonParser.FindStringSubmatch(str)
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
		_b[v] = field.LivingCell
	}
	for _, digit := range sDigits {
		v, e := strconv.Atoi(digit)
		if e != nil {
			panic(e)
		}
		_s[v] = field.LivingCell
	}

	return &Automaton{_b, _s}
}

// Do the rules allow survival for the given neighbour count?
func (ca *Automaton) Survival(neighbourCount int) field.Cell {
	return ca.S[neighbourCount]
}

// Do the rules allow giving birth for the given neighbour count?
func (ca *Automaton) Birth(neighbourCount int) field.Cell {
	return ca.B[neighbourCount]
}

// Next returns the state of the specified cell at the next time step.
func (ca *Automaton) NextCellState(f *field.Field, x, y int, cur field.Cell) field.Cell {

	neighbourCount := f.NeighbourCount(x, y)

	if cur.On() {
		return ca.Survival(neighbourCount)
	}

	return ca.Birth(neighbourCount)
}

// Apply a CA to the field
func (ca *Automaton) ApplyToField(f *field.Field) {
	f.MapAsync(ca.NextCellState)
}

func (ca *Automaton) String() string {
	var buf strings.Builder
	buf.WriteString("B")
	for i := 0; i <= 8; i++ {
		if ca.B[i].On() {
			buf.WriteString(strconv.Itoa(i))
		}
	}
	buf.WriteString("/S")
	for i := 0; i <= 8; i++ {
		if ca.S[i].On() {
			buf.WriteString(strconv.Itoa(i))
		}
	}
	return buf.String()
}
