package glox

type Parser struct {
	tokens  *[]*Token
	current int
}

func NewParser(tokens *[]*Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (s *Parser) ParseExpressionForTest() (Expr, error) {
	return s.expression()
}

func (s *Parser) Parse() ([]Stmt, error) {
	// program        → declaration* EOF ;
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
	for !s.check(TokenRightBrace) && !s.isAtEnd() {
		stmt, err := s.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	_, err := s.consume(TokenRightBrace, "Expect '}' after block.")
	if err != nil {
		return nil, err
	} else {
		return statements, nil
	}
}

// declaration    → classDecl
//                | funDecl
//                | varDecl
//                | statement ;
func (s *Parser) declaration() (Stmt, error) {
	if s.match(TokenClass) {
		return s.classDeclaration()
	}
	if s.match(TokenFun) {
		return s.function("function")
	}
	if s.match(TokenVar) {
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

// classDecl      → "class" IDENTIFIER ( "<" IDENTIFIER )?
//                  "{" function* "}" ;
func (s *Parser) classDeclaration() (Stmt, error) {
	className, err := s.consume(TokenIdentifier, "Expect class name.")
	if err != nil {
		return nil, err
	}
	var superclass *Variable
	if s.match(TokenLess) {
		_, err = s.consume(TokenIdentifier, "Expect superclass name.")
		if err != nil {
			return nil, err
		}
		superclass = NewVariable(s.previous())
	}
	_, err = s.consume(TokenLeftBrace, "Expect '{' before class body.")
	if err != nil {
		return nil, err
	}
	var methods []*Function
	for !s.check(TokenRightBrace) && !s.isAtEnd() {
		f, err := s.function("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, f)
	}
	_, err = s.consume(TokenRightBrace, "Expect '}' after class body.")
	if err != nil {
		return nil, err
	}
	return NewClass(className, superclass, &methods), nil
}

// funDecl        → "fun" function ;
// function       → IDENTIFIER "(" parameters? ")" block ;
func (s *Parser) function(kind string) (*Function, error) {
	funcName, err := s.consume(TokenIdentifier, "Expect "+kind+" name.")
	if err != nil {
		return nil, err
	}
	_, err = s.consume(TokenLeftParen, "Expect '(' after "+kind+" name.")
	if err != nil {
		return nil, err
	}
	var parameters []*Token
	if !s.check(TokenRightParen) {
		for {
			// if len(parameters) >= 255 {
			// 	return nil, NewParserError(s.peek(), "Can't have more than 255 arguments.")
			// }

			parameterName, err := s.consume(TokenIdentifier, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, parameterName)

			if !s.match(TokenComma) {
				break
			}
		}
	}
	_, err = s.consume(TokenRightParen, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = s.consume(TokenLeftBrace, "Expect '{' before "+kind+" body.")
	if err != nil {
		return nil, err
	}

	body, err := s.block()
	if err != nil {
		return nil, err
	}
	return NewFunction(funcName, &parameters, &body), nil
}

// parameters     → IDENTIFIER ( "," IDENTIFIER )* ;

// varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
func (s *Parser) varDeclaration() (Stmt, error) {
	name, err := s.consume(TokenIdentifier, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	var initializer Expr = nil
	if s.match(TokenEqual) {
		initializer, err = s.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = s.consume(TokenSemicolon, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}
	return NewVar(name, initializer), nil
}

// statement      → exprStmt
//                | forStmt
//                | ifStmt
//                | printStmt
//                | returnStmt
//                | whileStmt
//                | block ;
func (s *Parser) statement() (Stmt, error) {
	if s.match(TokenFor) {
		return s.forStatement()
	}
	if s.match(TokenIf) {
		return s.ifStatement()
	}
	if s.match(TokenPrint) {
		return s.printStatement()
	}
	if s.match(TokenReturn) {
		return s.returnStatement()
	}
	if s.match(TokenWhile) {
		return s.whileStatement()
	}
	if s.match(TokenLeftBrace) {
		stmts, err := s.block()
		if err != nil {
			return nil, err
		}
		return NewBlock(&stmts), nil
	}
	return s.expressionStatement()
}

// exprStmt       → expression ";" ;
func (s *Parser) expressionStatement() (Stmt, error) {
	expr, err := s.expression()
	if err != nil {
		return nil, err
	}
	_, err = s.consume(TokenSemicolon, "Expect ';' after expression.")
	return NewExpression(expr), nil
}

// printStmt      → "print" expression ";" ;
func (s *Parser) printStatement() (Stmt, error) {
	value, err := s.expression()
	if err != nil {
		return nil, err
	}
	_, err = s.consume(TokenSemicolon, "Expect ';' after value.")
	return NewPrint(value), nil
}

// forStmt        → "for" "(" ( varDecl | exprStmt | ";" )
//                 expression? ";"
//                 expression? ")" statement ;
func (s *Parser) forStatement() (Stmt, error) {
	_, err := s.consume(TokenLeftParen, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer Stmt
	if s.match(TokenSemicolon) {
		initializer = nil
	} else if s.match(TokenVar) {
		initializer, err = s.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = s.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition Expr
	if !s.check(TokenSemicolon) {
		condition, err = s.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = s.consume(TokenSemicolon, "Expect ';' after for condition.")
	if err != nil {
		return nil, err
	}

	var increment Expr
	if !s.check(TokenRightParen) {
		increment, err = s.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = s.consume(TokenRightParen, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	body, err := s.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		var blockList []Stmt
		blockList = append(blockList, body)
		blockList = append(blockList, NewExpression(increment))
		body = NewBlock(&blockList)
	}

	if condition == nil {
		condition = NewLiteral(true)
	}
	body = NewWhile(condition, body)

	if initializer != nil {
		var blockList []Stmt
		blockList = append(blockList, initializer)
		blockList = append(blockList, body)
		body = NewBlock(&blockList)
	}

	return body, nil
}

// ifStmt         → "if" "(" expression ")" statement
//               ( "else" statement )? ;
func (s *Parser) ifStatement() (Stmt, error) {
	_, err := s.consume(TokenLeftParen, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, err := s.expression()
	if err != nil {
		return nil, err
	}
	_, err = s.consume(TokenRightParen, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}
	thenBranch, err := s.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch Stmt
	if s.match(TokenElse) {
		elseBranch, err = s.statement()
		if err != nil {
			return nil, err
		}
	}
	return NewIf(condition, thenBranch, elseBranch), nil
}

// whileStmt      → "while" "(" expression ")" statement ;
func (s *Parser) whileStatement() (Stmt, error) {
	_, err := s.consume(TokenLeftParen, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}
	condition, err := s.expression()
	if err != nil {
		return nil, err
	}
	_, err = s.consume(TokenRightParen, "Expect ')' after while condition.")
	if err != nil {
		return nil, err
	}
	body, err := s.statement()
	if err != nil {
		return nil, err
	}
	return NewWhile(condition, body), nil
}

// returnStmt     → "return" expression? ";" ;
func (s *Parser) returnStatement() (Stmt, error) {
	keyword := s.previous()
	var value Expr
	var err error
	if !s.check(TokenSemicolon) {
		value, err = s.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = s.consume(TokenSemicolon, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}
	return NewReturn(keyword, value), nil
}

// =====

// expression     → assignment ;
func (s *Parser) expression() (Expr, error) {
	// return s.equality()
	return s.assignment()
}

// assignment     → ( call "." )? IDENTIFIER "=" assignment
//                | logic_or ;
func (s *Parser) assignment() (Expr, error) {
	expr, err := s.or()
	if err != nil {
		return nil, err
	}
	if s.match(TokenEqual) {
		equals := s.previous()
		value, err := s.assignment()
		if err != nil {
			return nil, err
		}
		if v, ok := expr.(*Variable); ok {
			name := v.name
			return NewAssign(name, value), nil
		} else if get, ok := expr.(*Get); ok {
			return NewSet(get.object, get.name, value), nil
		}
		return nil, NewParserError(equals, "Invalid assignment target.")
	}
	return expr, nil
}

// logic_or       → logic_and ( "or" logic_and )* ;
func (s *Parser) or() (Expr, error) {
	expr, err := s.and()
	if err != nil {
		return nil, err
	}
	for s.match(TokenOr) {
		operator := s.previous()
		right, err := s.and()
		if err != nil {
			return nil, err
		}
		expr = NewLogical(expr, operator, right)
	}
	return expr, nil
}

// logic_and      → equality ( "and" equality )* ;
func (s *Parser) and() (Expr, error) {
	expr, err := s.equality()
	if err != nil {
		return nil, err
	}
	for s.match(TokenAnd) {
		operator := s.previous()
		right, err := s.equality()
		if err != nil {
			return nil, err
		}
		expr = NewLogical(expr, operator, right)
	}
	return expr, nil
}

// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
func (s *Parser) equality() (Expr, error) {
	expr, err := s.comparison()
	if err != nil {
		return nil, err
	}
	for s.match(TokenBangEqual, TokenEqualEqual) {
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
	for s.match(TokenGreater, TokenGreaterEqual, TokenLess, TokenLessEqual) {
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
	for s.match(TokenMinus, TokenPlus) {
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
	for s.match(TokenSlash, TokenStar) {
		operator := s.previous()
		right, err := s.unary()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}
	return expr, nil
}

// unary          → ( "!" | "-" ) unary | call ;
func (s *Parser) unary() (Expr, error) {
	if s.match(TokenBang, TokenMinus) {
		operator := s.previous()
		right, err := s.unary()
		if err != nil {
			return nil, err
		}
		return NewUnary(operator, right), nil
	}
	return s.call()
}

// call           → primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
func (s *Parser) call() (Expr, error) {
	expr, err := s.primary()
	if err != nil {
		return nil, err
	}
	for {
		if s.match(TokenLeftParen) {
			expr, err = s.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if s.match(TokenDot) {
			name, err := s.consume(TokenIdentifier, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = NewGet(expr, name)
		} else {
			break
		}
	}
	return expr, nil
}

func (s *Parser) finishCall(callee Expr) (Expr, error) {
	var arguments []Expr
	if !s.check(TokenRightParen) {
		for {
			// if len(arguments) >= 255 {
			// 	return nil, NewParserError(s.peek(), "Can't have more than 255 arguments.")
			// }

			expr, err := s.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, expr)

			if !s.match(TokenComma) {
				break
			}
		}
	}
	paren, err := s.consume(TokenRightParen, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}
	return NewCall(callee, paren, &arguments), nil
}

// arguments      → expression ( "," expression )* ;

// primary        → "true" | "false" | "nil" | "this"
//                | NUMBER | STRING | IDENTIFIER | "(" expression ")"
//                | "super" "." IDENTIFIER ;
func (s *Parser) primary() (Expr, error) {
	if s.match(TokenNumber, TokenString) {
		return NewLiteral(s.previous().literal), nil
	}
	if s.match(TokenSuper) {
		keyword := s.previous()
		_, err := s.consume(TokenDot, "Expect '.' after 'super'.")
		if err != nil {
			return nil, err
		}
		method, err := s.consume(TokenIdentifier, "Expect superclass method name.")
		if err != nil {
			return nil, err
		}
		return NewSuper(keyword, method), nil
	}
	if s.match(TokenThis) {
		return NewThis(s.previous()), nil
	}
	if s.match(TokenTrue) {
		return NewLiteral(true), nil
	}
	if s.match(TokenFalse) {
		return NewLiteral(false), nil
	}
	if s.match(TokenNil) {
		return NewLiteral(nil), nil
	}
	if s.match(TokenIdentifier) {
		return NewVariable(s.previous()), nil
	}
	if s.match(TokenLeftParen) {
		expr, err := s.expression()
		if err != nil {
			return nil, err
		}
		_, err = s.consume(TokenRightParen, "Expect ')' after expression.")
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
	return s.peek().tokenType == TokenEof
}

func (s *Parser) peek() *Token {
	return (*s.tokens)[s.current]
}

func (s *Parser) previous() *Token {
	return (*s.tokens)[s.current-1]
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
		if s.previous().tokenType == TokenSemicolon {
			return
		}
		switch s.peek().tokenType {
		case TokenClass:
		case TokenFun:
		case TokenVar:
		case TokenFor:
		case TokenIf:
		case TokenWhile:
		case TokenPrint:
		case TokenReturn:
			return
		}
		s.advance()
	}
}
