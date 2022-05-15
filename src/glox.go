package glox

import "fmt"

func Run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()
	for _, token := range tokens {
		// fmt.Println(token.String())
		fmt.Print(token.Lexeme())
		fmt.Print(" | ")
	}
	fmt.Println()
}
