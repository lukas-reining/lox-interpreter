package expressions

import (
	"github.com/lukas-reining/lox/scanner"
)

type Binary[T any, Err error] struct {
	Expression[T, Err]

	Left     Expression[T, Err]
	Operator scanner.Token
	Right    Expression[T, Err]
}

func NewBinary[T any, Err error](left Expression[T, Err], operator scanner.Token, right Expression[T, Err]) *Binary[T, Err] {
	return &Binary[T, Err]{
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func (e *Binary[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitBinaryExpression(e)
}
