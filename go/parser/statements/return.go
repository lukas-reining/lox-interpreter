package statements

import (
	"github.com/lukas-reining/lox/parser/expressions"
	"github.com/lukas-reining/lox/scanner"
)

type Return[T any, Err error] struct {
	Statement[T, Err]

	Keyword scanner.Token
	Value   expressions.Expression[T, Err]
}

func NewReturn[T any, Err error](keyword scanner.Token, value expressions.Expression[T, Err]) *Return[T, Err] {
	return &Return[T, Err]{
		Keyword: keyword,
		Value:   value,
	}
}

func (e *Return[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitReturnStatement(e)
}
