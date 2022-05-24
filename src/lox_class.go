package glox

type LoxClass struct {
	name    string
	methods *map[string]*LoxFunction
}

func NewLoxClass(name string, methods *map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		name:    name,
		methods: methods,
	}
}

func (s *LoxClass) arity() int {
	initializer := s.findMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.arity()
}

func (s *LoxClass) call(interpreter *Interpreter, arguments *[]interface{}) (interface{}, error) {
	instance := NewLoxInstance(s)
	initializer := s.findMethod("init")
	if initializer != nil {
		method, err := initializer.bind(instance)
		if err != nil {
			return nil, err
		}
		_, err = method.call(interpreter, arguments)
		if err != nil {
			return nil, err
		}
	}
	return instance, nil
}

func (s *LoxClass) findMethod(name string) *LoxFunction {
	if v, ok := (*s.methods)[name]; ok {
		return v
	}
	return nil
}

func (s *LoxClass) String() string {
	return s.name
}
