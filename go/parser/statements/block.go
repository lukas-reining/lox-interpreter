package statements

type Block[T any, Err error] struct {
	Statement[T, Err]

	Statements []Statement[T, Err]
}

func NewBlock[T any, Err error](statements []Statement[T, Err]) *Block[T, Err] {
	return &Block[T, Err]{
		Statements: statements,
	}
}

func (e *Block[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitBlockStatement(e)
}
