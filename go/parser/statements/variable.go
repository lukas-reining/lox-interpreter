package statements

import (
	"github.com/lukas-reining/lox/parser/expressions"
	"github.com/lukas-reining/lox/scanner"
)

type Var[T any, Err error] struct {
	Statement[T, Err]

	Name        scanner.Token
	Initializer expressions.Expression[T, Err]
}

func NewVar[T any, Err error](name scanner.Token, initializer expressions.Expression[T, Err]) *Var[T, Err] {
	return &Var[T, Err]{
		Name: name, Initializer: initializer,
	}
}

func (e *Var[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitVarStatement(e)
}
