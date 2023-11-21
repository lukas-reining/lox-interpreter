package expressions

type Literal[T any, Err error] struct {
	Expression[T, Err]

	Literal any
}

func NewLiteral[T any, Err error](literal any) *Literal[T, Err] {
	return &Literal[T, Err]{
		Literal: literal,
	}
}

func (e *Literal[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitLiteralExpression(e)
}
