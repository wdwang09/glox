package glox

import "fmt"

// https://craftinginterpreters.com/scanning.html#error-handling
// static boolean hadError = false;
// https://craftinginterpreters.com/evaluating-expressions.html#reporting-runtime-errors
// static boolean hadRuntimeError = false;

type LineError struct {
	line    int
	message string
}

func (s *LineError) Error() string {
	return fmt.Sprintf("[line %v] Error: %v\n", s.line, s.message)
}

func NewLineError(line int, message string) *LineError {
	return &LineError{
		line:    line,
		message: message,
	}
}

// =====

type ParserError struct {
	token   *Token
	message string
}

func (s *ParserError) Error() string {
	where := "end"
	if s.token.tokenType != TokenEof {
		where = "\"" + s.token.String() + "\""
	}
	return fmt.Sprintf("[line %v] Error at %v: %v\n",
		s.token.line, where, s.message)
}

func NewParserError(token *Token, message string) *ParserError {
	return &ParserError{
		token:   token,
		message: message,
	}
}

// =====

type RuntimeError struct {
	token   *Token
	message string
}

func (s *RuntimeError) Error() string {
	return fmt.Sprintf("[line %v] RuntimeError at \"%v\": %v",
		s.token.line, s.token.String(), s.message)
}

func NewRuntimeError(token *Token, message string) *RuntimeError {
	return &RuntimeError{
		token:   token,
		message: message,
	}
}
