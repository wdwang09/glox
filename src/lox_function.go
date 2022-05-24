package glox

import "time"

type LoxCallable interface {
	arity() int
	call(interpreter *Interpreter, arguments *[]interface{}) (interface{}, error)
}

// =====

type LoxFunction struct {
	declaration   *Function
	closure       *Environment
	isInitializer bool
}

func NewLoxFunction(declaration *Function, closure *Environment, isInitializer bool) *LoxFunction {
	return &LoxFunction{
		declaration:   declaration,
		closure:       closure,
		isInitializer: isInitializer,
	}
}

func (s *LoxFunction) arity() int {
	return len(*s.declaration.params)
}

func (s *LoxFunction) call(interpreter *Interpreter, arguments *[]interface{}) (interface{}, error) {
	environment := NewEnvironment(s.closure)
	for i, param := range *s.declaration.params {
		err := environment.define(param.lexeme, (*arguments)[i])
		if err != nil {
			return nil, err
		}
	}
	err := interpreter.executeBlock(s.declaration.body, environment)
	if returnValue, ok := err.(*ReturnPseudoError); ok {
		if s.isInitializer {
			return s.closure.getAt(0, "this")
		}
		return returnValue.value, nil
	}
	if s.isInitializer {
		return s.closure.getAt(0, "this")
	}
	return nil, err
}

func (s *LoxFunction) bind(instance *LoxInstance) (*LoxFunction, error) {
	environment := NewEnvironment(s.closure)
	err := environment.define("this", instance)
	if err != nil {
		return nil, err
	}
	return NewLoxFunction(s.declaration, environment, s.isInitializer), nil
}

func (s *LoxFunction) String() string {
	return "<Function " + s.declaration.name.lexeme + ">"
}

// =====

type clockLoxFunction struct{}

func NewClockLoxFunction() *clockLoxFunction {
	return &clockLoxFunction{}
}

func (s *clockLoxFunction) arity() int {
	return 0
}

func (s *clockLoxFunction) call(_ *Interpreter, _ *[]interface{}) (interface{}, error) {
	return float64(time.Now().UnixMilli()) / 1000.0, nil
}

func (s *clockLoxFunction) String() string {
	return "<Function clock>"
}
