package glox

import (
	"fmt"
)

func Run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()
	for _, token := range tokens {
		// fmt.Println(token.String())
		fmt.Print(token.Lexeme())
		fmt.Print(" | ")
	}
	fmt.Println()
	parser := NewParser(tokens)
	expression, err := parser.Parse()
	if err != nil {
		// fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	astPrinter := NewAstPrinter()
	fmt.Println(astPrinter.Print(expression))
}
