package interpeter

type LoxValue = any

type Stringifyable interface {
	ToString() string
}

type ReturnValue struct {
	Value LoxValue
}

func NewReturnValue(value LoxValue) *ReturnValue {
	return &ReturnValue{
		Value: value,
	}
}
