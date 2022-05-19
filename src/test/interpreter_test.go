package glox

import (
	"glox/src"
	"testing"
)

var expressionValue map[string]string

func initExpressionValueTestCase() {
	expressionValue = map[string]string{
		// primary
		"true":      "true",
		"12345":     "12345",
		"\"hello\"": "hello",
		"nil":       "nil",
		"(1)":       "1",

		// unary
		"!true":  "false",
		"!!true": "true",
		"-1":     "-1",
		"--1":    "1",
		"-(-1)":  "1",

		// factor
		"-1 * -1": "1",

		// term
		"1 + 2 * 3":       "7",
		"1 * 2 + 3 * 4":   "14",
		"1 * (2 + 3) * 4": "20",

		// comparison
		"1 + 2 * 3 > 2 * 2": "true",

		// equality
		"1 > 2 == 3 > 4":   "true",
		"1 <= 2 != 3 >= 4": "true",
		"\"123\" == 123":   "false",
		"\"nil\" != nil":   "true",
		"nil == nil":       "true",
	}
}

func TestInterpreter(t *testing.T) {
	initExpressionValueTestCase()
	for code, valueExpectation := range expressionValue {
		scanner := glox.NewScanner(code)
		tokens, err := scanner.ScanTokens()
		if err != nil {
			t.Fatal(err.Error())
		}
		parser := glox.NewParser(tokens)
		expression, err := parser.ParseExpressionForTest()
		if err != nil {
			t.Fatal(err.Error())
		}
		interpreter := glox.NewInterpreter()
		value, err := interpreter.InterpretExpressionForTest(expression)
		if err != nil {
			t.Fatal(err.Error())
		}
		value = interpreter.Stringify(value)
		if value != valueExpectation {
			t.Fatalf("\nTestcase: %v\nOutput: %v\nExpect: %v",
				code, value, valueExpectation)
		}
	}
}
