package statements

import (
	"github.com/lukas-reining/lox/parser/expressions"
	"github.com/lukas-reining/lox/scanner"
)

type Class[T any, Err error] struct {
	Statement[T, Err]

	Name    scanner.Token
	Super   expressions.Variable[T, Err]
	Methods []Function[T, Err]
}

func NewClass[T any, Err error](name scanner.Token, methods []Function[T, Err]) *Class[T, Err] {
	return &Class[T, Err]{
		Name:    name,
		Methods: methods,
	}
}

func (e *Class[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitClassStatement(e)
}
