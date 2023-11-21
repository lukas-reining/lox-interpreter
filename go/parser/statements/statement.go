package statements

type Statement[T any, Err error] interface {
	Accept(visitor Visitor[T, Err]) (T, Err)
}
