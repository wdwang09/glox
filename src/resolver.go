package glox

type Resolver struct {
	interpreter     *Interpreter
	scopes          *scopeStack
	currentFunction FunctionType
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          NewScopeStack(),
		currentFunction: FNone,
	}
}

func (s *Resolver) beginScope() {
	s.scopes.Push()
}

func (s *Resolver) endScope() {
	s.scopes.Pop()
}

func (s *Resolver) declare(name *Token) error {
	if s.scopes.IsEmpty() {
		return nil
	}
	scope := s.scopes.Peek()
	if _, ok := (*scope)[name.lexeme]; ok {
		return NewParserError(name, "Already a variable with this name in this scope.")
	}
	(*scope)[name.lexeme] = false
	return nil
}

func (s *Resolver) define(name *Token) {
	if s.scopes.IsEmpty() {
		return
	}
	(*(s.scopes.Peek()))[name.lexeme] = true
}

func (s *Resolver) ResolveStatements(statements *[]Stmt) error {
	for _, stmt := range *statements {
		err := s.resolveStatement(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Resolver) resolveStatement(stmt Stmt) error {
	_, err := stmt.Accept(s)
	return err
}

func (s *Resolver) resolveExpression(expr Expr) error {
	_, err := expr.Accept(s)
	return err
}

func (s *Resolver) resolveFunction(function *Function, fType FunctionType) error {
	enclosingFunction := s.currentFunction
	s.currentFunction = fType
	s.beginScope()
	for _, param := range *function.params {
		err := s.declare(param)
		if err != nil {
			return err
		}
		s.define(param)
	}
	err := s.ResolveStatements(function.body)
	if err != nil {
		return nil
	}
	s.endScope()
	s.currentFunction = enclosingFunction
	return nil
}

// =====

func (s *Resolver) visitBlockStmt(stmt *Block) (interface{}, error) {
	s.beginScope()
	err := s.ResolveStatements(stmt.statements)
	if err != nil {
		return nil, err
	}
	s.endScope()
	return nil, nil
}

func (s *Resolver) visitClassStmt(stmt *Class) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Resolver) visitExpressionStmt(stmt *Expression) (interface{}, error) {
	return nil, s.resolveExpression(stmt.expression)
}

func (s *Resolver) visitFunctionStmt(stmt *Function) (interface{}, error) {
	err := s.declare(stmt.name)
	if err != nil {
		return nil, err
	}
	s.define(stmt.name)
	return nil, s.resolveFunction(stmt, FFunction)
}

func (s *Resolver) visitIfStmt(stmt *If) (interface{}, error) {
	err := s.resolveExpression(stmt.condition)
	if err != nil {
		return nil, err
	}
	err = s.resolveStatement(stmt.thenBranch)
	if err != nil {
		return nil, err
	}
	if stmt.elseBranch != nil {
		err = s.resolveStatement(stmt.elseBranch)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (s *Resolver) visitPrintStmt(stmt *Print) (interface{}, error) {
	return nil, s.resolveExpression(stmt.expression)
}

func (s *Resolver) visitReturnStmt(stmt *Return) (interface{}, error) {
	if s.currentFunction == FNone {
		return nil, NewParserError(stmt.keyword, "Can't return from top-level code.")
	}
	if stmt.value != nil {
		return nil, s.resolveExpression(stmt.value)
	}
	return nil, nil
}

func (s *Resolver) visitVarStmt(stmt *Var) (interface{}, error) {
	err := s.declare(stmt.name)
	if err != nil {
		return nil, err
	}
	if stmt.initializer != nil {
		err = s.resolveExpression(stmt.initializer)
		if err != nil {
			return nil, err
		}
	}
	s.define(stmt.name)
	return nil, nil
}

func (s *Resolver) visitWhileStmt(stmt *While) (interface{}, error) {
	err := s.resolveExpression(stmt.condition)
	if err != nil {
		return nil, err
	}
	return nil, s.resolveStatement(stmt.body)
}

// =====

func (s *Resolver) visitAssignExpr(expr *Assign) (interface{}, error) {
	err := s.resolveExpression(expr.value)
	if err != nil {
		return nil, err
	}
	s.resolveLocal(expr, expr.name)
	return nil, nil
}

func (s *Resolver) visitBinaryExpr(expr *Binary) (interface{}, error) {
	err := s.resolveExpression(expr.left)
	if err != nil {
		return nil, err
	}
	return nil, s.resolveExpression(expr.right)
}

func (s *Resolver) visitCallExpr(expr *Call) (interface{}, error) {
	err := s.resolveExpression(expr.callee)
	if err != nil {
		return nil, err
	}
	for _, argument := range *expr.arguments {
		err = s.resolveExpression(argument)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (s *Resolver) visitGetExpr(expr *Get) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Resolver) visitGroupingExpr(expr *Grouping) (interface{}, error) {
	return nil, s.resolveExpression(expr.expression)
}

func (s *Resolver) visitLiteralExpr(_ *Literal) (interface{}, error) {
	return nil, nil
}

func (s *Resolver) visitLogicalExpr(expr *Logical) (interface{}, error) {
	err := s.resolveExpression(expr.left)
	if err != nil {
		return nil, err
	}
	return nil, s.resolveExpression(expr.right)
}

func (s *Resolver) visitSetExpr(expr *Set) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Resolver) visitSuperExpr(expr *Super) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Resolver) visitThisExpr(expr *This) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Resolver) visitUnaryExpr(expr *Unary) (interface{}, error) {
	return nil, s.resolveExpression(expr.right)
}

func (s *Resolver) visitVariableExpr(expr *Variable) (interface{}, error) {
	if !s.scopes.IsEmpty() {
		if value, ok := (*(s.scopes.Peek()))[expr.name.lexeme]; !value && ok {
			return nil, NewParserError(expr.name, "Can't read local variable in its own initializer.")
		}
	}
	s.resolveLocal(expr, expr.name)
	return nil, nil
}

func (s *Resolver) resolveLocal(expr Expr, name *Token) {
	for i := s.scopes.Size() - 1; i >= 0; i-- {
		if _, ok := (*s.scopes.Get(i))[name.lexeme]; ok {
			s.interpreter.Resolve(expr, s.scopes.Size()-1-i)
			return
		}
	}
}

// =====

type scopeStack struct {
	scopes []map[string]bool
}

func NewScopeStack() *scopeStack {
	return &scopeStack{
		scopes: []map[string]bool{},
	}
}

func (s *scopeStack) Push() {
	s.scopes = append(s.scopes, make(map[string]bool))
}

func (s *scopeStack) Pop() {
	s.scopes = s.scopes[:len(s.scopes)-1]
}

func (s *scopeStack) Get(i int) *map[string]bool {
	return &s.scopes[i]
}

func (s *scopeStack) Peek() *map[string]bool {
	return s.Get(s.Size() - 1)
}

func (s *scopeStack) IsEmpty() bool {
	return len(s.scopes) == 0
}

func (s *scopeStack) Size() int {
	return len(s.scopes)
}

// =====

type FunctionType int

const (
	FNone FunctionType = iota
	FFunction
)
