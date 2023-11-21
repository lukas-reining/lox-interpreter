package expressions

import (
	"github.com/lukas-reining/lox/scanner"
)

type Assignment[T any, Err error] struct {
	Expression[T, Err]

	Name  scanner.Token
	Value Expression[T, Err]
}

func NewAssignment[T any, Err error](name scanner.Token, value Expression[T, Err]) *Assignment[T, Err] {
	return &Assignment[T, Err]{
		Name:  name,
		Value: value,
	}
}

func (e *Assignment[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitAssignmentExpression(e)
}
