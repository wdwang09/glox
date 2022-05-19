package glox

import (
	"bufio"
	"fmt"
	"os"
)

type Glox struct {
	interpreter *Interpreter
}

func NewGlox() *Glox {
	return &Glox{interpreter: NewInterpreter()}
}

func (s *Glox) RunFile(path string) int {
	fileData, err := os.ReadFile(path)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "[File]", err)
		return 1
	}
	return s.run(string(fileData))
}

func (s *Glox) RunPrompt() int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				return 0
			}
			_, _ = fmt.Fprintln(os.Stderr, "[Prompt]", err)
			return 1
		}
		// if line == "" {
		//	continue
		// }
		s.run(line)
	}
}

func (s *Glox) run(source string) int {
	// Scanner
	scanner := NewScanner(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "[Scanner]", err.Error())
		return 1
	}
	fmt.Print("[Scanner] ")
	for _, token := range tokens {
		// fmt.Println(token.String())
		fmt.Print(token.Lexeme())
		fmt.Print(" | ")
	}
	fmt.Println()

	// Parser
	parser := NewParser(tokens)
	statements, err := parser.Parse()
	// expression, err := parser.Parse()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "[Parser]", err.Error())
		return 1
	}

	// AST Printer
	// astPrinter := NewAstPrinter()
	// fmt.Println(astPrinter.Print(expression))

	// Interpreter
	// err = s.interpreter.Interpret(statements)
	value, err := s.interpreter.Interpret(statements)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "[Interpreter]", err.Error())
		return 1
	}
	if value != nil {
		fmt.Println(s.interpreter.Stringify(value))
	}
	return 0
}
