package statements

import (
	"github.com/lukas-reining/lox/parser/expressions"
)

type Expression[T any, Err error] struct {
	Statement[T, Err]

	Exp expressions.Expression[T, Err]
}

func NewExpression[T any, Err error](exp expressions.Expression[T, Err]) *Expression[T, Err] {
	return &Expression[T, Err]{
		Exp: exp,
	}
}

func (e *Expression[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitExpressionStatement(e)
}
