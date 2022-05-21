package glox

import "strconv"

type Scanner struct {
	tokenMap *map[string]TokenType
	source   string
	tokens   []*Token
	start    int
	current  int
	line     int
}

func NewScanner(tokenMap *map[string]TokenType, source string) *Scanner {
	return &Scanner{
		tokenMap: tokenMap,
		source:   source,
		tokens:   []*Token{},
		start:    0,
		current:  0,
		line:     1,
	}
}

func (s *Scanner) ScanTokens() ([]*Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			return nil, err
		}
	}
	s.tokens = append(s.tokens, NewToken(TokenEof, "", nil, s.line))
	return s.tokens, nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() error {
	ch := s.advance()
	switch ch {
	case '(':
		s.addToken(TokenLeftParen)
	case ')':
		s.addToken(TokenRightParen)
	case '{':
		s.addToken(TokenLeftBrace)
	case '}':
		s.addToken(TokenRightBrace)
	case ',':
		s.addToken(TokenComma)
	case '.':
		s.addToken(TokenDot)
	case '-':
		s.addToken(TokenMinus)
	case '+':
		s.addToken(TokenPlus)
	case ';':
		s.addToken(TokenSemicolon)
	case '*':
		s.addToken(TokenStar)

	case '!':
		if s.match('=') {
			s.addToken(TokenBangEqual)
		} else {
			s.addToken(TokenBang)
		}
	case '=':
		if s.match('=') {
			s.addToken(TokenEqualEqual)
		} else {
			s.addToken(TokenEqual)
		}
	case '<':
		if s.match('=') {
			s.addToken(TokenLessEqual)
		} else {
			s.addToken(TokenLess)
		}
	case '>':
		if s.match('=') {
			s.addToken(TokenGreaterEqual)
		} else {
			s.addToken(TokenGreater)
		}

	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(TokenSlash)
		}

	case ' ', '\r', '\t':

	case '\n':
		s.line++

	case '"':
		err := s.string()
		if err != nil {
			return err
		}

	default:
		if isDigit(ch) {
			s.number()
		} else if isAlpha(ch) {
			s.identifier()
		} else {
			return NewLineError(s.line, "Unexpected character.")
		}
	}
	return nil
}

func (s *Scanner) advance() byte {
	ch := s.source[s.current]
	s.current++
	return ch
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenLiteral(tokenType, nil)
}

func (s *Scanner) addTokenLiteral(tokenType TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, text, literal, s.line))
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() || s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return NewLineError(s.line, "Unexpected character.")
	}

	s.advance()
	s.addTokenLiteral(TokenString, s.source[s.start+1:s.current-1])
	return nil
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	num, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)
	s.addTokenLiteral(TokenNumber, num)
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch == '_')
}

func iSAlphaNumeric(ch byte) bool {
	return isDigit(ch) || isAlpha(ch)
}

func (s *Scanner) identifier() {
	for iSAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, exist := (*s.tokenMap)[text]
	if !exist {
		tokenType = TokenIdentifier
	}
	s.addToken(tokenType)
}
