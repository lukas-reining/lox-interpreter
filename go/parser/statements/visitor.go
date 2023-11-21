package statements

type Visitor[T any, Err error] interface {
	VisitPrintStatement(exp *Print[T, Err]) (T, Err)
	VisitExpressionStatement(exp *Expression[T, Err]) (T, Err)
	VisitVarStatement(exp *Var[T, Err]) (T, Err)
	VisitBlockStatement(exp *Block[T, Err]) (T, Err)
	VisitIfStatement(exp *If[T, Err]) (T, Err)
	VisitWhileStatement(exp *While[T, Err]) (T, Err)
	VisitFunctionStatement(exp *Function[T, Err]) (T, Err)
	VisitReturnStatement(exp *Return[T, Err]) (T, Err)
	VisitClassStatement(exp *Class[T, Err]) (T, Err)
}
