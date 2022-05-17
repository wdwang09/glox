package glox

import (
	"glox/src"
	"testing"
)

var expressionAst map[string]string

func initExpressionAstTestCase() {
	expressionAst = map[string]string{
		// primary
		"true":      "true",
		"12345":     "12345",
		"\"hello\"": "hello",
		"nil":       "nil",
		"(1)":       "(group 1)",

		// unary
		"!true":  "(! true)",
		"!!true": "(! (! true))",
		"-1":     "(- 1)",
		"--1":    "(- (- 1))",
		"-(-1)":  "(- (group (- 1)))",

		// factor
		"-1 * -1": "(* (- 1) (- 1))",

		// term
		"1 + 2 * 3":       "(+ 1 (* 2 3))",
		"1 * 2 + 3 * 4":   "(+ (* 1 2) (* 3 4))",
		"1 * (2 + 3) * 4": "(* (* 1 (group (+ 2 3))) 4)",

		// comparison
		"1 + 2 * 3 > 2 * 2": "(> (+ 1 (* 2 3)) (* 2 2))",

		// equality
		"1 > 2 == 3 > 4":   "(== (> 1 2) (> 3 4))",
		"1 <= 2 != 3 >= 4": "(!= (<= 1 2) (>= 3 4))",
	}
}

func TestParser(t *testing.T) {
	initExpressionAstTestCase()
	for code, astExpectation := range expressionAst {
		scanner := glox.NewScanner(code)
		tokens, err := scanner.ScanTokens()
		if err != nil {
			t.Fatal(err.Error())
		}
		parser := glox.NewParser(tokens)
		expression, err := parser.Parse()
		if err != nil {
			t.Fatal(err.Error())
		}
		astPrinter := glox.NewAstPrinter()
		astOutput, err := astPrinter.Print(expression)
		if err != nil {
			t.Fatal(err.Error())
		}
		if astOutput != astExpectation {
			t.Fatalf("\nOutput: %v\nExpect: %v", astOutput, astExpectation)
		}
	}
}
