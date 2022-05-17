package glox

import "fmt"

// https://craftinginterpreters.com/scanning.html#error-handling
// static boolean hadError = false;
// https://craftinginterpreters.com/evaluating-expressions.html#reporting-runtime-errors
// static boolean hadRuntimeError = false;

type lineError struct {
	line    int
	message string
}

func (s *lineError) Error() string {
	return fmt.Sprintf("[line %v] Error: %v\n", s.line, s.message)
}

func NewLineError(line int, message string) *lineError {
	s := new(lineError)
	s.line = line
	s.message = message
	return s
}

// =====

type parserError struct {
	token   *Token
	message string
}

func (s *parserError) Error() string {
	where := "end"
	if s.token.tokenType != EOF {
		where = "\"" + s.token.String() + "\""
	}
	return fmt.Sprintf("[line %v] Error at %v: %v\n",
		s.token.line, where, s.message)
}

func NewParserError(token *Token, message string) *parserError {
	s := new(parserError)
	s.token = token
	s.message = message
	return s
}

// =====

type runtimeError struct {
	token   *Token
	message string
}

func (s *runtimeError) Error() string {
	return fmt.Sprintf("[line %v] RuntimeError at \"%v\": %v",
		s.token.line, s.token.String(), s.message)
}

func NewRuntimeError(token *Token, message string) *runtimeError {
	s := new(runtimeError)
	s.token = token
	s.message = message
	return s
}
