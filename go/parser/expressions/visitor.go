package expressions

type Visitor[T any, Err error] interface {
	VisitGroupingExpression(exp *Grouping[T, Err]) (T, Err)
	VisitBinaryExpression(exp *Binary[T, Err]) (T, Err)
	VisitUnaryExpression(exp *Unary[T, Err]) (T, Err)
	VisitLiteralExpression(exp *Literal[T, Err]) (T, Err)
	VisitVariableExpression(exp *Variable[T, Err]) (T, Err)
	VisitAssignmentExpression(exp *Assignment[T, Err]) (T, Err)
	VisitLogicalExpression(exp *Logical[T, Err]) (T, Err)
	VisitCallExpression(exp *Call[T, Err]) (T, Err)
	VisitGetExpression(exp *Get[T, Err]) (T, Err)
	VisitSetExpression(exp *Set[T, Err]) (T, Err)
	VisitThisExpression(exp *This[T, Err]) (T, Err)
}
