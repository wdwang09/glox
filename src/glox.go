package glox

import (
	"fmt"
	"os"
)

func Run(source string) {
	// Scanner
	scanner := NewScanner(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	for _, token := range tokens {
		// fmt.Println(token.String())
		fmt.Print(token.Lexeme())
		fmt.Print(" | ")
	}
	fmt.Println()

	// Parser
	parser := NewParser(tokens)
	expression, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	// AST Printer
	astPrinter := NewAstPrinter()
	fmt.Println(astPrinter.Print(expression))

	// Interpreter
	interpreter := NewInterpreter()
	value, err := interpreter.Interpret(expression)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	fmt.Println(interpreter.Stringify(value))
}
