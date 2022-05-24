package glox

type Resolver struct {
	interpreter     *Interpreter
	scopes          *scopeStack
	currentFunction FunctionType
	currentClass    ClassType
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          NewScopeStack(),
		currentFunction: FNone,
		currentClass:    CNone,
	}
}

func (s *Resolver) beginScope() {
	s.scopes.push()
}

func (s *Resolver) endScope() {
	s.scopes.pop()
}

func (s *Resolver) declare(name *Token) error {
	if s.scopes.isEmpty() {
		return nil
	}
	scope := s.scopes.peek()
	if _, ok := (*scope)[name.lexeme]; ok {
		return NewResolverError(name, "Already a variable with this name in this scope.")
	}
	(*scope)[name.lexeme] = false
	return nil
}

func (s *Resolver) define(name *Token) {
	if s.scopes.isEmpty() {
		return
	}
	(*(s.scopes.peek()))[name.lexeme] = true
}

func (s *Resolver) resolveStatements(statements *[]Stmt) error {
	for _, stmt := range *statements {
		err := s.resolveStatement(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Resolver) resolveStatement(stmt Stmt) error {
	_, err := stmt.accept(s)
	return err
}

func (s *Resolver) resolveExpression(expr Expr) error {
	_, err := expr.accept(s)
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
	err := s.resolveStatements(function.body)
	if err != nil {
		return err
	}
	s.endScope()
	s.currentFunction = enclosingFunction
	return nil
}

// =====

func (s *Resolver) visitBlockStmt(stmt *Block) (interface{}, error) {
	s.beginScope()
	err := s.resolveStatements(stmt.statements)
	if err != nil {
		return nil, err
	}
	s.endScope()
	return nil, nil
}

func (s *Resolver) visitClassStmt(stmt *Class) (interface{}, error) {
	enclosingClass := s.currentClass
	s.currentClass = CClass
	err := s.declare(stmt.name)
	if err != nil {
		return nil, err
	}
	s.define(stmt.name)
	if stmt.superclass != nil && stmt.name.lexeme == stmt.superclass.name.lexeme {
		return nil, NewResolverError(stmt.superclass.name, "A class can't inherit from itself.")
	}
	if stmt.superclass != nil {
		s.currentClass = CSubclass
		err = s.resolveExpression(stmt.superclass)
		if err != nil {
			return nil, err
		}
		s.beginScope()
		(*s.scopes.peek())["super"] = true
	}
	s.beginScope()
	(*s.scopes.peek())["this"] = true
	for _, method := range *stmt.methods {
		declaration := FMethod
		if method.name.lexeme == "init" {
			declaration = FInitializer
		}
		err = s.resolveFunction(method, declaration)
		if err != nil {
			return nil, err
		}
	}
	s.endScope()
	if stmt.superclass != nil {
		s.endScope()
	}
	s.currentClass = enclosingClass
	return nil, nil
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
		return nil, NewResolverError(stmt.keyword, "Can't return from top-level code.")
	}
	if stmt.value != nil {
		if s.currentFunction == FInitializer {
			return nil, NewResolverError(stmt.keyword, "Can't return a value from an initializer.")
		}
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
	return nil, s.resolveExpression(expr.object)
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
	err := s.resolveExpression(expr.value)
	if err != nil {
		return nil, err
	}
	return nil, s.resolveExpression(expr.object)
}

func (s *Resolver) visitSuperExpr(expr *Super) (interface{}, error) {
	if s.currentClass == CNone {
		return nil, NewResolverError(expr.keyword, "Can't use 'super' outside of a class.")
	} else if s.currentClass != CSubclass {
		return nil, NewResolverError(expr.keyword, "Can't use 'super' in a class with no superclass.")
	}
	s.resolveLocal(expr, expr.keyword)
	return nil, nil
}

func (s *Resolver) visitThisExpr(expr *This) (interface{}, error) {
	if s.currentClass == CNone {
		return nil, NewResolverError(expr.keyword, "Can't use 'this' outside of a class.")
	}
	s.resolveLocal(expr, expr.keyword)
	return nil, nil
}

func (s *Resolver) visitUnaryExpr(expr *Unary) (interface{}, error) {
	return nil, s.resolveExpression(expr.right)
}

func (s *Resolver) visitVariableExpr(expr *Variable) (interface{}, error) {
	if !s.scopes.isEmpty() {
		if value, ok := (*(s.scopes.peek()))[expr.name.lexeme]; !value && ok {
			return nil, NewResolverError(expr.name, "Can't read local variable in its own initializer.")
		}
	}
	s.resolveLocal(expr, expr.name)
	return nil, nil
}

func (s *Resolver) resolveLocal(expr Expr, name *Token) {
	for i := s.scopes.size() - 1; i >= 0; i-- {
		if _, ok := (*s.scopes.get(i))[name.lexeme]; ok {
			s.interpreter.resolve(expr, s.scopes.size()-1-i)
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

func (s *scopeStack) push() {
	s.scopes = append(s.scopes, make(map[string]bool))
}

func (s *scopeStack) pop() {
	s.scopes = s.scopes[:len(s.scopes)-1]
}

func (s *scopeStack) get(i int) *map[string]bool {
	return &s.scopes[i]
}

func (s *scopeStack) peek() *map[string]bool {
	return s.get(s.size() - 1)
}

func (s *scopeStack) isEmpty() bool {
	return len(s.scopes) == 0
}

func (s *scopeStack) size() int {
	return len(s.scopes)
}

// =====

type FunctionType int

const (
	FNone FunctionType = iota
	FFunction
	FInitializer
	FMethod
)

type ClassType int

const (
	CNone ClassType = iota
	CClass
	CSubclass
)
