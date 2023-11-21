package statements

import "github.com/lukas-reining/lox/parser/expressions"

type If[T any, Err error] struct {
	Statement[T, Err]

	Condition  expressions.Expression[T, Err]
	IfBranch   Statement[T, Err]
	ElseBranch Statement[T, Err]
}

func NewIf[T any, Err error](
	condition expressions.Expression[T, Err],
	ifBranch Statement[T, Err],
	elseBranch Statement[T, Err],
) *If[T, Err] {
	return &If[T, Err]{
		Condition: condition, IfBranch: ifBranch, ElseBranch: elseBranch,
	}
}

func (e *If[T, Err]) Accept(visitor Visitor[T, Err]) (T, Err) {
	return visitor.VisitIfStatement(e)
}
