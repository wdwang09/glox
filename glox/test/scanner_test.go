package glox

import (
	"fmt"
	"glox/glox"
	"testing"
)

var loxCode map[string]int

func initTestCase() {
	loxCode = map[string]int{
		"print \"hello\";":          4,
		"1 + 1 == 2;":               7,
		"var a = 3; a == 2;\n":      10,
		"var a = 3; a == 2;\r\n":    10,
		"var ord = 3; 1 or 2;\r\n":  10,
		"var ord = 3;\n1 or 2;\r\n": 10,
	}
}

func TestScanner(t *testing.T) {
	initTestCase()
	for code, num := range loxCode {
		scanner := glox.NewScanner(code)
		tokens := scanner.ScanTokens()
		if len(tokens) != num {
			fmt.Println("====================================================")
			fmt.Println(code)
			for _, token := range tokens {
				fmt.Print(token.Lexeme())
				fmt.Print(" | ")
			}
			fmt.Println()
			fmt.Println("====================================================")
			t.Fatalf(`the tokens' length of code "%v" should be %d, get %d`, code, num, len(tokens))
		}
	}
}
