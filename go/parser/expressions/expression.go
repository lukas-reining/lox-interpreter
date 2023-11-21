package expressions

type Expression[T any, Err error] interface {
	Accept(visitor Visitor[T, Err]) (T, Err)
}
