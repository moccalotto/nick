package field

import (
	"strconv"
	"strings"
	"regexp"
)

// Rules for evolving a field.
type Rules struct {
	// lookup tables
	B [9]bool	// if B[X] is true, it means that a cell will be born if it has X neighbours
	S [9]bool	// if S[X] is true, it means that a cell will survive if it has X neighbours
}

// Create new Rules
func MakeRules(b, s [9]bool) *Rules {
	return &Rules{b, s}
}

func NewRules(str string) *Rules {
	_s := [9]bool{}
	_b := [9]bool{}

	re := regexp.MustCompile("^B([0-9]*)/S([0-9]*)$")
	matches := re.FindStringSubmatch(str)
	if (len(matches) != 3) {
		// TODO: correct error handling.
		panic("Bad Rule-String")
	}

	bDigits := strings.Split(matches[1], "")
	sDigits := strings.Split(matches[2], "")

	for _, digit := range(bDigits) {
		v,e := strconv.Atoi(digit)
		if (e != nil) {
			panic(e)
		}
		_b[v] = true
	}
	for _, digit := range(sDigits) {
		v,e := strconv.Atoi(digit)
		if (e != nil) {
			panic(e)
		}
		_s[v] = true
	}

	return &Rules{_b, _s}
}

// Do the rules allow survival for the given neighbour count?
func (r *Rules) Survival(neighbourCount int) bool {
	return r.S[neighbourCount]
}

// Do the rules allow giving birth for the given neighbour count?
func (r *Rules) Birth(neighbourCount int) bool {
	return r.B[neighbourCount]
}

// Calculate the next state of a given cell
func (r *Rules) NextCellState(currentlyAlive bool, neighbourCount int) bool {
	if currentlyAlive {
		return r.Survival(neighbourCount)
	} else {
		return r.Birth(neighbourCount)
	}
}

func (r *Rules) String() string {
	var buf strings.Builder
	buf.WriteString("B")
	for i := 0; i <= 8; i++ {
		if r.B[i] {
			buf.WriteString(strconv.Itoa(i))
		}
	}
	buf.WriteString("/S")
	for i := 0; i <= 8; i++ {
		if r.S[i] {
			buf.WriteString(strconv.Itoa(i))
		}
	}
	return buf.String()
}
