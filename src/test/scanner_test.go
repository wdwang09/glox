package glox

import (
	"fmt"
	"glox/src"
	"testing"
)

var loxCodeNum map[string]int

func initLoxCodeNumTestCase() {
	loxCodeNum = map[string]int{
		"print \"hello\";":         4,
		"1 + 1 == 2;":              7,
		"var a = 3; a == 2;\n":     10,
		"var a = 3; a == 2;\r\n":   10,
		"var ord = 3; 1 or 2;\r\n": 10,
		"var ord = 3;\n1 or 2;\n":  10,
	}
}

func TestScanner(t *testing.T) {
	initLoxCodeNumTestCase()
	tokenMap := glox.NewTokenMap()
	for code, num := range loxCodeNum {
		scanner := glox.NewScanner(tokenMap, code)
		tokens, err := scanner.ScanTokens()
		if err != nil {
			t.Fatal(err.Error())
		}
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
