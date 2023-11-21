package expressions

type Grouping[T any, Err error] struct {
	Expression[T, Err]
	Exp Expression[T, Err]
}

func NewGrouping[T any, Err error](expression Expression[T, Err]) *Grouping[T, Err] {
	return &Grouping[T, Err]{
		Exp: expression,
	}
}

func (e *Grouping[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitGroupingExpression(e)
}
