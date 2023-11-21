package expressions

import (
	"github.com/lukas-reining/lox/scanner"
)

type Unary[T any, Err error] struct {
	Expression[T, Err]

	Operator scanner.Token
	Right    Expression[T, Err]
}

func NewUnary[T any, Err error](operator scanner.Token, right Expression[T, Err]) *Unary[T, Err] {
	return &Unary[T, Err]{
		Operator: operator,
		Right:    right,
	}
}

func (e *Unary[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitUnaryExpression(e)
}
