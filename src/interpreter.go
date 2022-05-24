package glox

import (
	"fmt"
	"reflect"
)

type Interpreter struct {
	globals     *Environment
	environment *Environment
	locals      map[Expr]int
}

func NewInterpreter() *Interpreter {
	environment := NewEnvironment(nil)
	_ = environment.Define("clock", NewClockLoxFunction())
	return &Interpreter{
		globals:     environment,
		environment: environment,
		locals:      map[Expr]int{},
	}
}

func (s *Interpreter) InterpretExpressionForTest(expr Expr) (interface{}, error) {
	value, err := s.evaluate(expr)
	if err != nil {
		return nil, err
	}
	return value, err
}

func (s *Interpreter) Interpret(statements *[]Stmt) (interface{}, error) {
	var value interface{}
	for _, stmt := range *statements {
		var err error
		value, err = s.execute(stmt)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

func (s *Interpreter) execute(stmt Stmt) (interface{}, error) {
	return stmt.Accept(s)
}

func (s *Interpreter) Stringify(obj interface{}) string {
	if obj == nil {
		return "nil"
	}
	if isFloat64(obj) {
		return fmt.Sprintf("%v", obj.(float64))
	}
	if isBool(obj) {
		return fmt.Sprintf("%v", obj.(bool))
	}
	return fmt.Sprint(obj)
}

// =====

func (s *Interpreter) visitBlockStmt(stmt *Block) (interface{}, error) {
	err := s.ExecuteBlock(stmt.statements, NewEnvironment(s.environment))
	return nil, err
}

func (s *Interpreter) visitClassStmt(stmt *Class) (interface{}, error) {
	err := s.environment.Define(stmt.name.lexeme, nil)
	if err != nil {
		return nil, err
	}
	methods := make(map[string]*LoxFunction)
	for _, method := range *stmt.methods {
		methods[method.name.lexeme] = NewLoxFunction(method, s.environment, method.name.lexeme == "init")
	}
	class := NewLoxClass(stmt.name.lexeme, &methods)
	err = s.environment.Assign(stmt.name, class)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *Interpreter) visitExpressionStmt(stmt *Expression) (interface{}, error) {
	return s.evaluate(stmt.expression)
}

func (s *Interpreter) visitFunctionStmt(stmt *Function) (interface{}, error) {
	function := NewLoxFunction(stmt, s.environment, false)
	return nil, s.environment.Define(stmt.name.lexeme, function)
}

func (s *Interpreter) visitIfStmt(stmt *If) (interface{}, error) {
	condition, err := s.evaluate(stmt.condition)
	if err != nil {
		return nil, err
	}
	if s.isTruthy(condition) {
		return s.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		return s.execute(stmt.elseBranch)
	}
	return nil, nil
}

func (s *Interpreter) visitPrintStmt(stmt *Print) (interface{}, error) {
	value, err := s.evaluate(stmt.expression)
	if err != nil {
		return nil, err
	}
	fmt.Println("[Print]", s.Stringify(value))
	return nil, nil
}

func (s *Interpreter) visitReturnStmt(stmt *Return) (interface{}, error) {
	var value interface{}
	var err error
	if stmt.value != nil {
		value, err = s.evaluate(stmt.value)
		if err != nil {
			return nil, err
		}
	}
	return nil, NewReturnPseudoError(value)
}

func (s *Interpreter) visitVarStmt(stmt *Var) (interface{}, error) {
	var value interface{} = nil
	if stmt.initializer != nil {
		var err error
		value, err = s.evaluate(stmt.initializer)
		if err != nil {
			return nil, err
		}
	}
	return nil, s.environment.Define(stmt.name.lexeme, value)
}

func (s *Interpreter) visitWhileStmt(stmt *While) (interface{}, error) {
	for {
		condition, err := s.evaluate(stmt.condition)
		if err != nil {
			return nil, err
		}
		if s.isTruthy(condition) {
			_, err = s.execute(stmt.body)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, nil
		}
	}
}

// =====

func (s *Interpreter) ExecuteBlock(statements *[]Stmt, environment *Environment) (err error) {
	previous := s.environment
	s.environment = environment
	for _, statement := range *statements {
		_, err = s.execute(statement)
		if err != nil {
			break
		}
	}
	s.environment = previous
	return err
}

// =====

func (s *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(s)
}

func (s *Interpreter) visitAssignExpr(expr *Assign) (interface{}, error) {
	value, err := s.evaluate(expr.value)
	if err != nil {
		return nil, err
	}

	// err = s.environment.Assign(expr.name, value)
	// if err != nil {
	// 	return nil, err
	// }
	// return value, nil

	if distance, ok := s.locals[expr]; ok {
		err = s.environment.AssignAt(distance, expr.name, value)
	} else {
		err = s.globals.Assign(expr.name, value)
	}
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (s *Interpreter) visitBinaryExpr(expr *Binary) (interface{}, error) {
	left, err := s.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	right, err := s.evaluate(expr.right)
	if err != nil {
		return nil, err
	}
	switch expr.operator.tokenType {
	case TokenGreater:
		err = s.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case TokenGreaterEqual:
		err = s.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case TokenLess:
		err = s.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case TokenLessEqual:
		err = s.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case TokenBangEqual:
		return !s.isEqual(left, right), nil
	case TokenEqualEqual:
		return s.isEqual(left, right), nil
	case TokenMinus:
		err = s.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case TokenSlash:
		err = s.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case TokenStar:
		err = s.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case TokenPlus:
		if isFloat64(left) && isFloat64(right) {
			return left.(float64) + right.(float64), nil
		}
		if isString(left) && isString(right) {
			return left.(string) + right.(string), nil
		}
		return nil, NewRuntimeError(expr.operator, "Operands must be two numbers or two strings.")
	}
	return nil, nil
}

func (s *Interpreter) visitCallExpr(expr *Call) (interface{}, error) {
	callee, err := s.evaluate(expr.callee)
	if err != nil {
		return nil, err
	}
	var arguments []interface{}
	for _, argument := range *expr.arguments {
		expr, err := s.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, expr)
	}
	if function, ok := callee.(LoxCallable); ok {
		if len(arguments) != function.arity() {
			return nil, NewRuntimeError(expr.paren,
				fmt.Sprintf("Expected %v arguments but got %v.", function.arity(), len(arguments)))
		}
		return function.call(s, &arguments)
	} else {
		return nil, NewRuntimeError(expr.paren, "Can only call functions and classes.")
	}
}

func (s *Interpreter) visitGetExpr(expr *Get) (interface{}, error) {
	obj, err := s.evaluate(expr.object)
	if err != nil {
		return nil, err
	}
	if v, ok := obj.(*LoxInstance); ok {
		return v.Get(expr.name)
	}
	return nil, NewRuntimeError(expr.name, "Only instances have properties.")
}

func (s *Interpreter) visitGroupingExpr(expr *Grouping) (interface{}, error) {
	return s.evaluate(expr.expression)
}

func (s *Interpreter) visitLiteralExpr(expr *Literal) (interface{}, error) {
	return expr.value, nil
}

func (s *Interpreter) visitLogicalExpr(expr *Logical) (interface{}, error) {
	left, err := s.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	if expr.operator.tokenType == TokenOr {
		if s.isTruthy(left) {
			return left, nil
		}
	} else if !s.isTruthy(left) { // AND
		return left, nil
	}
	return s.evaluate(expr.right)
}

func (s *Interpreter) visitSetExpr(expr *Set) (interface{}, error) {
	obj, err := s.evaluate(expr.object)
	if err != nil {
		return nil, err
	}
	if o, ok := obj.(*LoxInstance); !ok {
		return nil, NewRuntimeError(expr.name, "Only instances have fields.")
	} else {
		value, err := s.evaluate(expr.value)
		if err != nil {
			return nil, err
		}
		o.Set(expr.name, value)
		return value, nil
	}
}

func (s *Interpreter) visitSuperExpr(expr *Super) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Interpreter) visitThisExpr(expr *This) (interface{}, error) {
	return s.lookUpVariable(expr.keyword, expr)
}

func (s *Interpreter) visitUnaryExpr(expr *Unary) (interface{}, error) {
	right, err := s.evaluate(expr.right)
	if err != nil {
		return nil, err
	}
	switch expr.operator.tokenType {
	case TokenBang:
		return !s.isTruthy(right), nil
	case TokenMinus:
		err = s.checkNumberOperands(expr.operator, right)
		if err != nil {
			return nil, err
		}
		return -right.(float64), nil
	}
	return nil, nil
}

func (s *Interpreter) visitVariableExpr(expr *Variable) (interface{}, error) {
	// return s.environment.Get(expr.name)
	return s.lookUpVariable(expr.name, expr)
}

// ====

func (s *Interpreter) Resolve(expr Expr, depth int) {
	s.locals[expr] = depth
}

func (s *Interpreter) lookUpVariable(name *Token, expr Expr) (interface{}, error) {
	if distance, ok := s.locals[expr]; ok {
		return s.environment.GetAt(distance, name.lexeme)
	} else {
		return s.globals.Get(name)
	}
}

// =====

func (s *Interpreter) isTruthy(obj interface{}) bool {
	if obj == nil {
		return false
	}
	if isBool(obj) {
		return obj.(bool)
	}
	return true
}

func (s *Interpreter) isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func (s *Interpreter) checkNumberOperands(operator *Token, operands ...interface{}) error {
	for _, operand := range operands {
		if !isFloat64(operand) {
			return NewRuntimeError(operator, "Operand must be a number.")
		}
	}
	return nil
}

func isFloat64(obj interface{}) bool {
	return reflect.ValueOf(obj).Kind() == reflect.Float64
}

func isBool(obj interface{}) bool {
	return reflect.ValueOf(obj).Kind() == reflect.Bool
}

func isString(obj interface{}) bool {
	return reflect.ValueOf(obj).Kind() == reflect.String
}
