package glox

type TokenType int

const (
	// Single-character tokens.

	TokenLeftParen TokenType = iota
	TokenRightParen
	TokenLeftBrace
	TokenRightBrace
	TokenComma
	TokenDot
	TokenMinus
	TokenPlus
	TokenSemicolon
	TokenSlash
	TokenStar

	// One or two character tokens.

	TokenBang
	TokenBangEqual
	TokenEqual
	TokenEqualEqual
	TokenGreater
	TokenGreaterEqual
	TokenLess
	TokenLessEqual

	// Literals.

	TokenIdentifier
	TokenString
	TokenNumber

	// Keywords.

	TokenAnd
	TokenClass
	TokenElse
	TokenFalse
	TokenFun
	TokenFor
	TokenIf
	TokenNil
	TokenOr
	TokenPrint
	TokenReturn
	TokenSuper
	TokenThis
	TokenTrue
	TokenVar
	TokenWhile

	TokenEof
)

// =====

func NewTokenMap() *map[string]TokenType {
	tokenMap := map[string]TokenType{
		"and":    TokenAnd,
		"class":  TokenClass,
		"else":   TokenElse,
		"false":  TokenFalse,
		"for":    TokenFor,
		"fun":    TokenFun,
		"if":     TokenIf,
		"nil":    TokenNil,
		"or":     TokenOr,
		"print":  TokenPrint,
		"return": TokenReturn,
		"super":  TokenSuper,
		"this":   TokenThis,
		"true":   TokenTrue,
		"var":    TokenVar,
		"while":  TokenWhile,
	}
	return &tokenMap
}
