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
	return "Return pseudo-error."
}
