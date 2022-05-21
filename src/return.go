package glox

type ReturnPseudoError struct {
	value interface{}
}

func NewReturnPseudoError(value interface{}) *ReturnPseudoError {
	return &ReturnPseudoError{
		value: value,
	}
}

func (s *ReturnPseudoError) Error() string {
	return "If you see this in console, it means that one return occurs outside a function."
}
