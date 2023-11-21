package statements

import "github.com/lukas-reining/lox/parser/expressions"

type While[T any, Err error] struct {
	Statement[T, Err]

	Condition expressions.Expression[T, Err]
	Body      Statement[T, Err]
}

func NewWhile[T any, Err error](
	condition expressions.Expression[T, Err],
	body Statement[T, Err],
) *While[T, Err] {
	return &While[T, Err]{
		Condition: condition, Body: body,
	}
}

func (e *While[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitWhileStatement(e)
}
