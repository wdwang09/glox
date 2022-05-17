package glox

import (
	"fmt"
	"reflect"
)

type Interpreter struct{}

func NewInterpreter() *Interpreter {
	s := new(Interpreter)
	return s
}

func (s *Interpreter) Interpret(expr Expr) (interface{}, error) {
	value, err := s.evaluate(expr)
	if err != nil {
		return nil, err
	}
	return value, err
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

func (s *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(s)
}

func (s *Interpreter) visitAssignExpr(assign *Assign) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Interpreter) visitBinaryExpr(binary *Binary) (interface{}, error) {
	left, err := s.evaluate(binary.left)
	if err != nil {
		return nil, err
	}
	right, err := s.evaluate(binary.right)
	if err != nil {
		return nil, err
	}
	switch binary.operator.tokenType {
	case GREATER:
		err = s.checkNumberOperands(binary.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case GREATER_EQUAL:
		err = s.checkNumberOperands(binary.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case LESS:
		err = s.checkNumberOperands(binary.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		err = s.checkNumberOperands(binary.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case BANG_EQUAL:
		return !s.isEqual(left, right), nil
	case EQUAL_EQUAL:
		return s.isEqual(left, right), nil
	case MINUS:
		err = s.checkNumberOperands(binary.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case SLASH:
		err = s.checkNumberOperands(binary.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case STAR:
		err = s.checkNumberOperands(binary.operator, left, right)
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
		return nil, NewRuntimeError(binary.operator, "Operands must be two numbers or two strings.")
	}
	return nil, nil
}

func (s *Interpreter) visitCallExpr(call *Call) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Interpreter) visitGetExpr(get *Get) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Interpreter) visitGroupingExpr(grouping *Grouping) (interface{}, error) {
	return s.evaluate(grouping.expression)
}

func (s *Interpreter) visitLiteralExpr(literal *Literal) (interface{}, error) {
	return literal.value, nil
}

func (s *Interpreter) visitLogicalExpr(logical *Logical) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Interpreter) visitSetExpr(set *Set) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Interpreter) visitSuperExpr(super *Super) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Interpreter) visitThisExpr(this *This) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Interpreter) visitUnaryExpr(unary *Unary) (interface{}, error) {
	right, err := s.evaluate(unary.right)
	if err != nil {
		return nil, err
	}
	switch unary.operator.tokenType {
	case BANG:
		return !s.isTruthy(right), nil
	case MINUS:
		err = s.checkNumberOperands(unary.operator, right)
		if err != nil {
			return nil, err
		}
		return -right.(float64), nil
	}
	return nil, nil
}

func (s *Interpreter) visitVariableExpr(variable *Variable) (interface{}, error) {
	// TODO implement me
	panic("implement me")
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
