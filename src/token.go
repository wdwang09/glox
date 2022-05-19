package glox

import "fmt"

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      int
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{
		tokenType: tokenType,
		lexeme:    lexeme,
		literal:   literal,
		line:      line,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%v %v %v", t.tokenType, t.lexeme, t.literal)
}

func (t *Token) Lexeme() string {
	return t.lexeme
}
