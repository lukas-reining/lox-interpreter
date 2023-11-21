package expressions

import (
	"github.com/lukas-reining/lox/scanner"
)

type Logical[T any, Err error] struct {
	Expression[T, Err]

	Left     Expression[T, Err]
	Operator scanner.Token
	Right    Expression[T, Err]
}

func NewLogical[T any, Err error](left Expression[T, Err], operator scanner.Token, right Expression[T, Err]) *Logical[T, Err] {
	return &Logical[T, Err]{
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func (e *Logical[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitLogicalExpression(e)
}
