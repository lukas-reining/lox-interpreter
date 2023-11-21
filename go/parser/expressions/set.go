package expressions

import (
	"github.com/lukas-reining/lox/scanner"
)

type Set[T any, Err error] struct {
	Expression[T, Err]

	Object Expression[T, Err]
	Name   scanner.Token
	Value  Expression[T, Err]
}

func NewSet[T any, Err error](object Expression[T, Err], name scanner.Token, value Expression[T, Err]) *Set[T, Err] {
	return &Set[T, Err]{
		Object: object,
		Name:   name,
		Value:  value,
	}
}

func (e *Set[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitSetExpression(e)
}
