package machine

import (
	"github.com/moccalotto/nick/utils"
	"math"
	"math/rand"
	"strconv"
)

func (m *Machine) Calculator() *utils.Calculator {

	if m.calculator == nil {
		m.initCalculator()
	}

	return m.calculator
}

func (m *Machine) initCalculator() {
	m.calculator = utils.NewCalculator(

		func(varName string) float64 {
			s := m.MustGetVar(varName)

			res, err := strconv.ParseFloat(s, 64)

			if err != nil {
				m.Throw("Variable %s is not a number", varName)
			}

			return res
		},

		func(callName string, args []float64) float64 {
			switch callName {
			case "abs":
				return math.Abs(args[0])
			case "acos":
				return math.Acos(args[0])
			case "asin":
				return math.Asin(args[0])
			case "atan":
				return math.Atan(args[0])
			case "ceil":
				return math.Ceil(args[0])
			case "cos":
				return math.Cos(args[0])
			case "deg2rad":
				return args[0] * 2 * math.Pi / 360.0
			case "exp":
				return math.Exp(args[0])
			case "floor":
				return math.Floor(args[0])
			case "log":
				return math.Log(args[0])
			case "max":
				if args[0] > args[1] {
					return args[0]
				}
				return args[1]
			case "min":
				if args[0] < args[1] {
					return args[0]
				}
				return args[1]
			case "mod":
				return math.Mod(args[0], args[1])
			case "pi":
				return math.Pi
			case "pow":
				return math.Pow(args[0], args[1])
			case "rad2deg":
				return args[0] * 360.0 / (2 * math.Pi)
			case "rand":
				return rand.Float64()
			case "randint":
				max := int(args[0])
				if max <= 0 {
					return 0
				}
				return float64(rand.Intn(max))
			case "round":
				return math.Round(args[0])
			case "sin":
				return math.Sin(args[0])
			case "tan":
				return math.Tan(args[0])
			default:
				m.Throw("Unknown function: %s", callName)
			}

			panic("This code should never be reached!")
		},
	)
}
