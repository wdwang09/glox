package glox

import (
	"glox/src"
	"testing"
)

var testCaseExpr []glox.Expr
var testCaseOutput []string

func initAstTestCase() {
	var expression glox.Expr
	var output string

	expression = glox.NewBinary(
		glox.NewUnary(glox.NewToken(glox.MINUS, "-", nil, 1), glox.NewLiteral(123)),
		glox.NewToken(glox.STAR, "*", nil, 1),
		glox.NewGrouping(glox.NewLiteral(45.67)))
	output = "(* (- 123) (group 45.67))"

	testCaseExpr = append(testCaseExpr, expression)
	testCaseOutput = append(testCaseOutput, output)

	expression = glox.NewBinary(
		glox.NewBinary(
			glox.NewLiteral(123),
			glox.NewToken(glox.PLUS, "+", nil, 1),
			glox.NewLiteral(45.67)),
		glox.NewToken(glox.SLASH, "/", nil, 1),
		glox.NewLiteral(8.9))
	output = "(/ (+ 123 45.67) 8.9)"

	testCaseExpr = append(testCaseExpr, expression)
	testCaseOutput = append(testCaseOutput, output)
}

func TestPrintAST(t *testing.T) {
	initAstTestCase()
	astPrinter := glox.NewAstPrinter()
	for i, expr := range testCaseExpr {
		astOutput := astPrinter.Print(expr)
		if astOutput != testCaseOutput[i] {
			t.Fatalf("\nOutput: %v\nExpect: %v", astOutput, testCaseOutput[i])
		}
	}
}
