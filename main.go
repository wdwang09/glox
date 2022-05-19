package main

import (
	"fmt"
	"glox/src"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		_, _ = fmt.Fprintln(os.Stderr, "[Main] Usage: glox [script]")
		os.Exit(64)
	}

	loxInterpreter := glox.NewGlox()

	if len(os.Args) == 2 {
		code := loxInterpreter.RunFile(os.Args[1])
		if code != 0 {
			_, _ = fmt.Fprintln(os.Stderr, "[Main] Failed when running file", os.Args[1])
			os.Exit(code)
		}
	} else {
		code := loxInterpreter.RunPrompt()
		if code != 0 {
			os.Exit(code)
		}
	}
}
