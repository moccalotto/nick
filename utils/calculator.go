package utils

import (
	"fmt"
	"strconv"
	"unicode"
)

const EOF = rune(0)

// Calculator - a state machine for calculating simple mathematical expressions
type Calculator struct {
	stmt     string      // the statement to calculate.
	slen     int         // length (in runes) of the statement
	idx      int         // the current index in the statement.
	getVar   VarAccessor // a function with which to fetch variables.
	execCall CallFunc    // a function with which to execute calls (for instance sin() cos(), foo() )
}

// VarAccessor is the callback used to provide a value when a variable is requested.
type VarAccessor func(varName string) float64

// CallFunc is called when the calculator needs to call a function
type CallFunc func(callName string, args []float64) float64

// NewCalculator returns a Calculator instance
func NewCalculator(getVar VarAccessor, execCall CallFunc) *Calculator {
	return &Calculator{
		getVar:   getVar,
		execCall: execCall,
	}
}

// Eval evaluates a simple mathematical statement
func (calc *Calculator) Eval(stmt string) float64 {
	calc.idx = 0
	calc.slen = len(stmt)
	calc.stmt = stmt

	calc.skipWhite()

	res := calc.evalExpression()

	if calc.idx != calc.slen {
		panic(fmt.Errorf("Syntax error. Unexpected: %c", calc.peekRune()))
	}

	return res
}

// BNF: expression = term { addop term }
func (calc *Calculator) evalExpression() float64 {

	res := calc.evalTerm()

	for calc.nextIsAddOp() {
		op := calc.evalRune()
		operand := calc.evalTerm()

		switch op {
		case '+':
			res += operand
		case '-':
			res -= operand
		}
	}

	calc.skipWhite()

	return res
}

// BNF: term = factor { mulop factor }
func (calc *Calculator) evalTerm() float64 {
	res := calc.evalFactor()

	for calc.nextIsMulOp() {
		op := calc.evalRune()
		operand := calc.evalFactor()

		switch op {
		case '*':
			res *= operand
		case '/':
			res /= operand
		}
	}

	calc.skipWhite()

	return res
}

// BNF:
// factor = "(" expression ")"
//        | call
//        | var
//        | number
func (calc *Calculator) evalFactor() (res float64) {
	nextRune := calc.peekRune()

	if nextRune == '(' {
		// "(" expression ")"
		calc.evalRune()
		res = calc.evalExpression()
		calc.expectRune(')')

		return res
	}

	if nextRune == '$' {
		// "$" identifier
		calc.next() // forward without skipping whitespace

		return calc.getVar(calc.evalIdentifier())
	}

	if unicode.IsLetter(nextRune) {
		// identifier "(" [ expression { "," expression } ] ")"
		return calc.evalCall()
	}

	return calc.evalNumber()
}

// BNF:
// number          =  ["-"] pointfloat | exponentfloat
// pointfloat      =  [intpart] fraction | intpart "."
// exponentfloat   =  (intpart | pointfloat) exponent
// intpart         =  digit+
// fraction        =  "." digit+
// exponent        =  ("e" | "E") ["+" | "-"] digit+
func (calc *Calculator) evalNumber() float64 {
	str := ""
	if calc.peekRune() == '-' {
		calc.next()
		str = "-"
	}

	// read an intpart
	str += calc.readDigits()

	if calc.peekRune() == '.' {
		calc.next()
		str += "." + calc.readDigits()
	}

	if exp := calc.peekRune(); exp == 'e' || exp == 'E' {
		str += string(exp)
		calc.next()
	} else {
		calc.skipWhite()

		res, _ := strconv.ParseFloat(str, 64)

		return res
	}

	if sign := calc.peekRune(); sign == '+' || sign == '-' {
		str += string(sign)
		calc.next()
	}

	str += calc.readDigits()

	calc.skipWhite()

	res, _ := strconv.ParseFloat(str, 64)

	return res
}

// Read an integer, forward the index, but do NOT skip whitespace
func (calc *Calculator) readDigits() string {
	str := ""
	for c := calc.peekRune(); unicode.IsDigit(c); c = calc.peekRune() {
		str += string(c)
		calc.next()
	}

	return str
}

// BNF: call = identifier "(" [ expression { "," expression } ] ")"
func (calc *Calculator) evalCall() float64 {
	name := calc.evalIdentifier()
	args := make([]float64, 0, 1)

	calc.expectRune('(')

	for nextRune := calc.peekRune(); nextRune != ')'; nextRune = calc.peekRune() {
		if len(args) > 0 {
			calc.expectRune(',')
		}
		args = append(args, calc.evalExpression())
	}

	calc.expectRune(')')

	return calc.execCall(name, args)
}

// BNF: identifier = letter { letter | digit }  // without whitespace!
func (calc *Calculator) evalIdentifier() string {
	if !unicode.IsLetter(calc.peekRune()) {
		panic(fmt.Errorf("Identier expected, but %c found", calc.peekRune()))
	}

	res := ""

	for calc.isIdentRune(calc.peekRune()) {
		res += string(calc.peekRune())
		calc.next()
	}

	calc.skipWhite()

	return res
}

func (calc *Calculator) isIdentRune(r rune) bool {
	return r == '_' ||
		unicode.IsLetter(r) ||
		unicode.IsDigit(r)
}

func (calc *Calculator) nextIsAddOp() bool {
	if c := calc.peekRune(); c == '-' || c == '+' {
		return true
	} else {
		return false
	}
}

func (calc *Calculator) nextIsMulOp() bool {
	if c := calc.peekRune(); c == '/' || c == '*' {
		return true
	} else {
		return false
	}
}

func (calc *Calculator) evalRune() rune {
	res := calc.peekRune()

	calc.next()

	calc.skipWhite()

	return res
}

func (calc *Calculator) peekRune() rune {

	if calc.idx >= calc.slen {
		return EOF
	}

	return rune(calc.stmt[calc.idx])
}

func (calc *Calculator) expectRune(expected rune) {
	if char := calc.evalRune(); char != expected {
		panic(fmt.Errorf("Expected '%c', but got '%c' instead", char, expected))
	}
}

func (calc *Calculator) skipWhite() {
	for unicode.IsSpace(calc.peekRune()) {
		calc.next()
	}
}

func (calc *Calculator) next() {
	calc.idx++

	if calc.idx > calc.slen {
		panic("Unexpected end of statement!")
	}
}
