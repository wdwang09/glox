package main

import (
	"bufio"
	"fmt"
	"glox/glox"
	"os"
)

func runFile(path string) int {
	fileData, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	glox.Run(string(fileData))
	return 0
}

func runPrompt() int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		// if line == "" {
		//	continue
		// }
		glox.Run(line)
	}
}

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "usage: glox [script]")
		os.Exit(64)
	}

	if len(os.Args) == 2 {
		code := runFile(os.Args[1])
		if code != 0 {
			os.Exit(code)
		}
	} else {
		code := runPrompt()
		if code != 0 {
			os.Exit(code)
		}
	}
}
