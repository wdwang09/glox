package glox

import "fmt"

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      int
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) *Token {
	t := new(Token)
	t.tokenType = tokenType
	t.lexeme = lexeme
	t.literal = literal
	t.line = line
	return t
}

func (t *Token) String() string {
	return fmt.Sprintf("%v %v %v", t.tokenType, t.lexeme, t.literal)
}

func (t *Token) Lexeme() string {
	return t.lexeme
}
