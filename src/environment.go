package glox

import "fmt"

type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]interface{}),
	}
}

func (s *Environment) define(name string, value interface{}) {
	s.values[name] = value
}

func (s *Environment) get(name *Token) (interface{}, error) {
	if obj, ok := s.values[name.lexeme]; ok {
		return obj, nil
	} else if s.enclosing != nil {
		return s.enclosing.get(name)
	}
	return nil, NewRuntimeError(name, fmt.Sprintf("Undefined varibale '%v'.", name.lexeme))
}

func (s *Environment) assign(name *Token, value interface{}) error {
	if _, ok := s.values[name.lexeme]; ok {
		s.values[name.lexeme] = value
		return nil
	} else if s.enclosing != nil {
		return s.enclosing.assign(name, value)
	}
	return NewRuntimeError(name, fmt.Sprintf("Undefined varibale '%v'.", name.lexeme))
}
