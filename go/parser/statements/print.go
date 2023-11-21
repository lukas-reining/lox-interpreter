package statements

import (
	"github.com/lukas-reining/lox/parser/expressions"
)

type Print[T any, Err error] struct {
	Statement[T, Err]

	Exp expressions.Expression[T, Err]
}

func NewPrintStatement[T any, Err error](exp expressions.Expression[T, Err]) *Print[T, Err] {
	return &Print[T, Err]{
		Exp: exp,
	}
}

func (e *Print[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitPrintStatement(e)
}
