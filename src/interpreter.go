package glox

import (
	"fmt"
	"reflect"
)

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: NewEnvironment(nil),
	}
}

func (s *Interpreter) InterpretExpressionForTest(expr Expr) (interface{}, error) {
	value, err := s.evaluate(expr)
	if err != nil {
		return nil, err
	}
	return value, err
}

func (s *Interpreter) Interpret(statements []Stmt) (interface{}, error) {
	var value interface{}
	for _, stmt := range statements {
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
		text := fmt.Sprintf("%v", obj.(float64))
		return text
	}
	if isBool(obj) {
		text := fmt.Sprintf("%v", obj.(bool))
		return text
	}
	return obj.(string)
}

// =====

func (s *Interpreter) visitBlockStmt(stmt *Block) (interface{}, error) {
	err := s.executeBlock(stmt.statements, NewEnvironment(s.environment))
	return nil, err
}

func (s *Interpreter) visitClassStmt(stmt *Class) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Interpreter) visitExpressionStmt(stmt *Expression) (interface{}, error) {
	return s.evaluate(stmt.expression)
}

func (s *Interpreter) visitFunctionStmt(stmt *Function) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Interpreter) visitIfStmt(stmt *If) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Interpreter) visitPrintStmt(stmt *Print) (interface{}, error) {
	value, err := s.evaluate(stmt.expression)
	if err != nil {
		return nil, err
	}
	fmt.Println("[Lox]", s.Stringify(value))
	return nil, nil
}

func (s *Interpreter) visitReturnStmt(stmt *Return) (interface{}, error) {
	// TODO implement me
	panic("implement me")
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
	s.environment.define(stmt.name.lexeme, value)
	return nil, nil
}

func (s *Interpreter) visitWhileStmt(stmt *While) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

// =====

func (s *Interpreter) executeBlock(statements []Stmt, environment *Environment) (err error) {
	previous := s.environment
	s.environment = environment
	for _, statement := range statements {
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
	err = s.environment.assign(expr.name, value)
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
	case GREATER:
		err = s.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case GREATER_EQUAL:
		err = s.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case LESS:
		err = s.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		err = s.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case BANG_EQUAL:
		return !s.isEqual(left, right), nil
	case EQUAL_EQUAL:
		return s.isEqual(left, right), nil
	case MINUS:
		err = s.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case SLASH:
		err = s.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case STAR:
		err = s.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case PLUS:
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
	// TODO implement me
	panic("implement me")
}

func (s *Interpreter) visitGetExpr(expr *Get) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Interpreter) visitGroupingExpr(expr *Grouping) (interface{}, error) {
	return s.evaluate(expr.expression)
}

func (s *Interpreter) visitLiteralExpr(expr *Literal) (interface{}, error) {
	return expr.value, nil
}

func (s *Interpreter) visitLogicalExpr(expr *Logical) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Interpreter) visitSetExpr(expr *Set) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Interpreter) visitSuperExpr(expr *Super) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Interpreter) visitThisExpr(expr *This) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Interpreter) visitUnaryExpr(expr *Unary) (interface{}, error) {
	right, err := s.evaluate(expr.right)
	if err != nil {
		return nil, err
	}
	switch expr.operator.tokenType {
	case BANG:
		return !s.isTruthy(right), nil
	case MINUS:
		err = s.checkNumberOperands(expr.operator, right)
		if err != nil {
			return nil, err
		}
		return -right.(float64), nil
	}
	return nil, nil
}

func (s *Interpreter) visitVariableExpr(expr *Variable) (interface{}, error) {
	return s.environment.get(expr.name)
}

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
