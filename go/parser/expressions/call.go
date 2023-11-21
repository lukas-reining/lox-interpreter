package expressions

import (
	"github.com/lukas-reining/lox/scanner"
)

type Call[T any, Err error] struct {
	Expression[T, Err]

	Callee      Expression[T, Err]
	Parenthesis scanner.Token
	Params      []Expression[T, Err]
}

func NewCall[T any, Err error](callee Expression[T, Err], parenthesis scanner.Token, params []Expression[T, Err]) *Call[T, Err] {
	return &Call[T, Err]{
		Callee:      callee,
		Parenthesis: parenthesis,
		Params:      params,
	}
}

func (e *Call[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitCallExpression(e)
}
