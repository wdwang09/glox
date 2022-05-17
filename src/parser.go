package glox

type Parser struct {
	tokens  []*Token
	current int
}

func NewParser(tokens []*Token) *Parser {
	s := new(Parser)
	s.tokens = tokens
	s.current = 0
	return s
}

func (s *Parser) Parse() (Expr, error) {
	expr, err := s.expression()
	return expr, err
}

// expression     → equality ;
func (s *Parser) expression() (Expr, error) {
	return s.equality()
}

// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
func (s *Parser) equality() (Expr, error) {
	expr, err := s.comparison()
	if err != nil {
		return nil, err
	}
	for s.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := s.previous()
		right, err := s.comparison()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}
	return expr, nil
}

// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (s *Parser) comparison() (Expr, error) {
	expr, err := s.term()
	if err != nil {
		return nil, err
	}
	for s.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := s.previous()
		right, err := s.term()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}
	return expr, nil
}

// term           → factor ( ( "-" | "+" ) factor )* ;
func (s *Parser) term() (Expr, error) {
	expr, err := s.factor()
	if err != nil {
		return nil, err
	}
	for s.match(MINUS, PLUS) {
		operator := s.previous()
		right, err := s.factor()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}
	return expr, nil
}

// factor         → unary ( ( "/" | "*" ) unary )* ;
func (s *Parser) factor() (Expr, error) {
	expr, err := s.unary()
	if err != nil {
		return nil, err
	}
	for s.match(SLASH, STAR) {
		operator := s.previous()
		right, err := s.unary()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}
	return expr, nil
}

// unary          → ( "!" | "-" ) unary
//                | primary ;
func (s *Parser) unary() (Expr, error) {
	if s.match(BANG, MINUS) {
		operator := s.previous()
		right, err := s.unary()
		if err != nil {
			return nil, err
		}
		return NewUnary(operator, right), nil
	}
	return s.primary()
}

// primary        → NUMBER | STRING | "true" | "false" | "nil"
//                | "(" expression ")" ;
func (s *Parser) primary() (Expr, error) {
	if s.match(NUMBER, STRING) {
		return NewLiteral(s.previous().literal), nil
	}
	if s.match(TRUE) {
		return NewLiteral(true), nil
	}
	if s.match(FALSE) {
		return NewLiteral(false), nil
	}
	if s.match(NIL) {
		return NewLiteral(nil), nil
	}
	if s.match(LEFT_PAREN) {
		expr, err := s.expression()
		if err != nil {
			return nil, err
		}
		_, err = s.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return NewGrouping(expr), nil
	}
	return nil, NewParserError(s.peek(), "Expect expression.")
}

func (s *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if s.check(tokenType) {
			s.advance()
			return true
		}
	}
	return false
}

func (s *Parser) check(tokenType TokenType) bool {
	if s.isAtEnd() {
		return false
	}
	return s.peek().tokenType == tokenType
}

func (s *Parser) advance() *Token {
	if !s.isAtEnd() {
		s.current++
	}
	return s.previous()
}

func (s *Parser) isAtEnd() bool {
	return s.peek().tokenType == EOF
}

func (s *Parser) peek() *Token {
	return s.tokens[s.current]
}

func (s *Parser) previous() *Token {
	return s.tokens[s.current-1]
}

func (s *Parser) consume(tokenType TokenType, message string) (*Token, error) {
	if s.check(tokenType) {
		return s.advance(), nil
	}
	return nil, NewParserError(s.peek(), message)
}

func (s *Parser) synchronize() {
	s.advance()
	for !s.isAtEnd() {
		if s.previous().tokenType == SEMICOLON {
			return
		}
		switch s.peek().tokenType {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}
		s.advance()
	}
}
