package expressions

import (
	"github.com/lukas-reining/lox/scanner"
)

type This[T any, Err error] struct {
	Expression[T, Err]

	Keyword scanner.Token
}

func NewThis[T any, Err error](keyword scanner.Token) *This[T, Err] {
	return &This[T, Err]{
		Keyword: keyword,
	}
}

func (e *This[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitThisExpression(e)
}
