package machine

import (
	"github.com/moccalotto/nick/utils"
	"math"
	"strconv"
)

func init() {
	InstructionHandlers["set"] = Set
}

// set $foo = [[this is a long string]]
func Set(m *Machine) {
	a := m.Arg(0)
	errStr := "Invalid use of 'set' instruction.\n" +
		"correct use is one of the following:\n" +
		"    set $foo = some-value       - sets $foo = 'some-value'\n" +
		"    set $bar between 1 and 20   - sets $bar to a random number in the range [1, 20]\n" +
		"    set $baz oneof a b c d      - sets $baz to either 'a', 'b', 'c' or 'd'"

	m.Assert(a.T == VarArg, errStr)

	switch m.ArgAsString(1) {
	case "=":
		m.Vars[a.StrVal] = m.ArgAsString(2)
	case "oneof":
		m.Vars[a.StrVal] = m.ArgAsString(m.Rng.Intn(m.ArgCount()-2) + 2)
		// do nothing
	case "between":
		m.Assert(m.ArgAsString(3) == "and", errStr)
		fmin, fmax := sortArgs(m.ArgAsFloat(2), m.ArgAsFloat(4))
		imin, imax := int(fmin), int(fmax)

		// no real difference between the two args.
		if fmax-fmin < 0.0000001 {
			m.Vars[a.StrVal] = strconv.FormatFloat(fmin, 'f', -1, 64)
		} else if float64(imin) == fmin && float64(imax) == fmax {
			delta := imax - imin
			m.Vars[a.StrVal] = strconv.Itoa(m.Rng.Intn(delta) + imin)
		} else {
			delta := fmax - fmin
			f := m.Rng.Float64()*delta + fmin
			m.Vars[a.StrVal] = strconv.FormatFloat(f, 'f', -1, 64)
		}
	case "calc":
		statement := ""
		for i := 2; i < m.ArgCount(); i++ {
			statement += m.ArgAsString(i)
			calc := utils.NewCalculator(
				func(varName string) float64 {
					s := m.MustGetVar(varName)

					res, err := strconv.ParseFloat(s, 64)

					if err != nil {
						m.Throw("Variable %s is not a number", varName)
					}

					return res
				},
				func(callName string, args []float64) float64 {
					return math.NaN()
				},
			)

			m.Vars[a.StrVal] = strconv.FormatFloat(
				calc.Eval(statement),
				'f',
				-1,
				64,
			)
		}

	default:
		m.Throw(errStr)
	}
}

func sortArgs(a, b float64) (float64, float64) {
	if a < b {
		return a, b
	} else {
		return b, a
	}
}
