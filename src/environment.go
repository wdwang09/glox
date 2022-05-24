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

func (s *Environment) Define(name string, value interface{}) error {
	s.values[name] = value
	return nil
}

func (s *Environment) Get(name *Token) (interface{}, error) {
	if obj, ok := s.values[name.lexeme]; ok {
		return obj, nil
	} else if s.enclosing != nil {
		return s.enclosing.Get(name)
	}
	return nil, NewRuntimeError(name, fmt.Sprintf("Undefined varibale '%v'.", name.lexeme))
}

func (s *Environment) GetAt(distance int, name string) (interface{}, error) {
	// if obj, ok := s.ancestor(distance).values[name.lexeme]; ok {
	// 	return obj, nil
	// }
	// return nil, NewRuntimeError(name, fmt.Sprintf("Undefined variable '%v'.", name.lexeme))
	return s.ancestor(distance).values[name], nil
}

func (s *Environment) ancestor(distance int) *Environment {
	environment := s
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
}

func (s *Environment) Assign(name *Token, value interface{}) error {
	if _, ok := s.values[name.lexeme]; ok {
		s.values[name.lexeme] = value
		return nil
	} else if s.enclosing != nil {
		return s.enclosing.Assign(name, value)
	}
	return NewRuntimeError(name, fmt.Sprintf("Undefined varibale '%v'.", name.lexeme))
}

func (s *Environment) AssignAt(distance int, name *Token, value interface{}) error {
	return s.ancestor(distance).Assign(name, value)
}
