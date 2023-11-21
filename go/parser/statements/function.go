package statements

import (
	"github.com/lukas-reining/lox/scanner"
)

type Function[T any, Err error] struct {
	Statement[T, Err]

	Name   scanner.Token
	Params []scanner.Token
	Body   []Statement[T, Err]
}

func NewFunction[T any, Err error](name scanner.Token, params []scanner.Token, body []Statement[T, Err]) *Function[T, Err] {
	return &Function[T, Err]{
		Name:   name,
		Params: params,
		Body:   body,
	}
}

func (e *Function[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitFunctionStatement(e)
}
