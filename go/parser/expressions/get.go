package expressions

import (
	"github.com/lukas-reining/lox/scanner"
)

type Get[T any, Err error] struct {
	Expression[T, Err]

	Object Expression[T, Err]
	Name   scanner.Token
}

func NewGet[T any, Err error](value Expression[T, Err], name scanner.Token) *Get[T, Err] {
	return &Get[T, Err]{
		Object: value,
		Name:   name,
	}
}

func (e *Get[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitGetExpression(e)
}
