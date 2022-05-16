package glox

import (
	"fmt"
	"os"
)

// https://craftinginterpreters.com/scanning.html#error-handling
// static boolean hadError = false;

func loxLineError(line int, message string) {
	report(line, "", message)
}

func loxTokenError(token *Token, message string) {
	if token.tokenType == EOF {
		report(token.line, "at end", message)
	} else {
		report(token.line, "at '"+token.String()+"'", message)
	}
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %v] Error %v: %v\n", line, where, message)
}
