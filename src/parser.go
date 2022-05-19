package glox

type Parser struct {
	tokens  []*Token
	current int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (s *Parser) ParseExpressionForTest() (Expr, error) {
	return s.expression()
}

func (s *Parser) Parse() ([]Stmt, error) {
	var statements []Stmt
	for !s.isAtEnd() {
		stmt, err := s.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	return statements, nil
}

// =====

func (s *Parser) block() ([]Stmt, error) {
	var statements []Stmt
	for !s.check(RIGHT_BRACE) && !s.isAtEnd() {
		stmt, err := s.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	_, err := s.consume(RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	} else {
		return statements, nil
	}
}

// declaration    → varDecl
//                | statement ;
func (s *Parser) declaration() (Stmt, error) {
	if s.match(VAR) {
		stmt, err := s.varDeclaration()
		// other methods: https://go.dev/blog/go1.13-errors
		if _, ok := err.(*ParserError); ok {
			s.synchronize()
		}
		if err != nil {
			return nil, err
		}
		return stmt, nil
	}
	stmt, err := s.statement()
	if _, ok := err.(*ParserError); ok {
		s.synchronize()
	}
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func (s *Parser) varDeclaration() (Stmt, error) {
	name, err := s.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	var initializer Expr = nil
	if s.match(EQUAL) {
		initializer, err = s.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = s.consume(SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}
	return NewVar(name, initializer), nil
}

// statement      → exprStmt
//                | printStmt ;
func (s *Parser) statement() (Stmt, error) {
	if s.match(PRINT) {
		return s.printStatement()
	}
	if s.match(LEFT_BRACE) {
		stmts, err := s.block()
		if err != nil {
			return nil, err
		}
		return NewBlock(stmts), nil
	}
	return s.expressionStatement()
}

func (s *Parser) expressionStatement() (Stmt, error) {
	expr, err := s.expression()
	if err != nil {
		return nil, err
	}
	_, err = s.consume(SEMICOLON, "Expect ';' after expression.")
	return NewExpression(expr), nil
}

func (s *Parser) printStatement() (Stmt, error) {
	value, err := s.expression()
	if err != nil {
		return nil, err
	}
	_, err = s.consume(SEMICOLON, "Expect ';' after value.")
	return NewPrint(value), nil
}

// =====

// expression     → assignment ;
func (s *Parser) expression() (Expr, error) {
	// return s.equality()
	return s.assignment()
}

// assignment     → IDENTIFIER "=" assignment
//                | equality ;
func (s *Parser) assignment() (Expr, error) {
	expr, err := s.equality()
	if err != nil {
		return nil, err
	}
	if s.match(EQUAL) {
		equals := s.previous()
		value, err := s.assignment()
		if err != nil {
			return nil, err
		}
		if v, ok := expr.(*Variable); ok {
			name := v.name
			return NewAssign(name, value), nil
		}
		return nil, NewParserError(equals, "Invalid assignment target.")
	}
	return expr, nil
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
	if s.match(IDENTIFIER) {
		return NewVariable(s.previous()), nil
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
