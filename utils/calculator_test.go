package utils

import (
	"testing"
)

type testScenario struct {
	stmt     string
	expected float64
}

func TestCalculator(t *testing.T) {
	calc := NewCalculator(
		func(varName string) float64 {
			return 777.0
		},
		func(callName string, args []float64) float64 {
			res := 666.0
			for _, a := range args {
				res += a
			}
			return res
		},
	)

	scenarios := []testScenario{
		{"3 * 13         + 5.0    ", 44.0},
		{"$var", 777.0},
		{"$var * 3", 2331},
		{"func()", 666.0},
		{"func() * 2", 2.0 * 666.0},
		{"func(1)", 667.0},
		{"func(2)", 668.0},
		{"func(1,2,3,4,5)", 681.0},
		{"2 * (func() - $var)", -222.0},
	}

	for _, s := range scenarios {
		result := calc.Eval(s.stmt)
		if result != s.expected {
			t.Errorf("Expected %s to be %f, but %f was returned", s.stmt, s.expected, result)
		}
	}
}
