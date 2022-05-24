package glox

// This code is generated by a Python script.

type Stmt interface {
	Accept(visitor stmtVisitor) (interface{}, error)
}

type stmtVisitor interface {
	visitBlockStmt(stmt *Block) (interface{}, error)
	visitClassStmt(stmt *Class) (interface{}, error)
	visitExpressionStmt(stmt *Expression) (interface{}, error)
	visitFunctionStmt(stmt *Function) (interface{}, error)
	visitIfStmt(stmt *If) (interface{}, error)
	visitPrintStmt(stmt *Print) (interface{}, error)
	visitReturnStmt(stmt *Return) (interface{}, error)
	visitVarStmt(stmt *Var) (interface{}, error)
	visitWhileStmt(stmt *While) (interface{}, error)
}

type Block struct {
	statements *[]Stmt
}

func NewBlock(statements *[]Stmt) *Block {
	stmt := new(Block)
	stmt.statements = statements
	return stmt
}

func (stmt *Block) Accept(visitor stmtVisitor) (interface{}, error) {
	return visitor.visitBlockStmt(stmt)
}

type Class struct {
	name       *Token
	superclass *Variable
	methods    *[]*Function
}

func NewClass(name *Token, superclass *Variable, methods *[]*Function) *Class {
	stmt := new(Class)
	stmt.name = name
	stmt.superclass = superclass
	stmt.methods = methods
	return stmt
}

func (stmt *Class) Accept(visitor stmtVisitor) (interface{}, error) {
	return visitor.visitClassStmt(stmt)
}

type Expression struct {
	expression Expr
}

func NewExpression(expression Expr) *Expression {
	stmt := new(Expression)
	stmt.expression = expression
	return stmt
}

func (stmt *Expression) Accept(visitor stmtVisitor) (interface{}, error) {
	return visitor.visitExpressionStmt(stmt)
}

type Function struct {
	name   *Token
	params *[]*Token
	body   *[]Stmt
}

func NewFunction(name *Token, params *[]*Token, body *[]Stmt) *Function {
	stmt := new(Function)
	stmt.name = name
	stmt.params = params
	stmt.body = body
	return stmt
}

func (stmt *Function) Accept(visitor stmtVisitor) (interface{}, error) {
	return visitor.visitFunctionStmt(stmt)
}

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func NewIf(condition Expr, thenBranch Stmt, elseBranch Stmt) *If {
	stmt := new(If)
	stmt.condition = condition
	stmt.thenBranch = thenBranch
	stmt.elseBranch = elseBranch
	return stmt
}

func (stmt *If) Accept(visitor stmtVisitor) (interface{}, error) {
	return visitor.visitIfStmt(stmt)
}

type Print struct {
	expression Expr
}

func NewPrint(expression Expr) *Print {
	stmt := new(Print)
	stmt.expression = expression
	return stmt
}

func (stmt *Print) Accept(visitor stmtVisitor) (interface{}, error) {
	return visitor.visitPrintStmt(stmt)
}

type Return struct {
	keyword *Token
	value   Expr
}

func NewReturn(keyword *Token, value Expr) *Return {
	stmt := new(Return)
	stmt.keyword = keyword
	stmt.value = value
	return stmt
}

func (stmt *Return) Accept(visitor stmtVisitor) (interface{}, error) {
	return visitor.visitReturnStmt(stmt)
}

type Var struct {
	name        *Token
	initializer Expr
}

func NewVar(name *Token, initializer Expr) *Var {
	stmt := new(Var)
	stmt.name = name
	stmt.initializer = initializer
	return stmt
}

func (stmt *Var) Accept(visitor stmtVisitor) (interface{}, error) {
	return visitor.visitVarStmt(stmt)
}

type While struct {
	condition Expr
	body      Stmt
}

func NewWhile(condition Expr, body Stmt) *While {
	stmt := new(While)
	stmt.condition = condition
	stmt.body = body
	return stmt
}

func (stmt *While) Accept(visitor stmtVisitor) (interface{}, error) {
	return visitor.visitWhileStmt(stmt)
}
