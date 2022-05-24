package glox

type LoxInstance struct {
	class  *LoxClass
	fields *map[string]interface{}
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		class:  class,
		fields: &map[string]interface{}{},
	}
}

func (s *LoxInstance) Get(name *Token) (interface{}, error) {
	if value, ok := (*s.fields)[name.lexeme]; ok {
		return value, nil
	}
	method := s.class.findMethod(name.lexeme)
	if method != nil {
		return method.bind(s)
	}
	return nil, NewRuntimeError(name, "Undefined property '"+name.lexeme+"'.")
}

func (s *LoxInstance) Set(name *Token, value interface{}) {
	(*s.fields)[name.lexeme] = value
}

func (s *LoxInstance) String() string {
	return s.class.name + " instance"
}
