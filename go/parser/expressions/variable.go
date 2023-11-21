package expressions

import (
	"github.com/lukas-reining/lox/scanner"
)

type Variable[T any, Err error] struct {
	Expression[T, Err]

	Name scanner.Token
}

func NewVariable[T any, Err error](name scanner.Token) *Variable[T, Err] {
	return &Variable[T, Err]{
		Name: name,
	}
}

func (e *Variable[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitVariableExpression(e)
}
