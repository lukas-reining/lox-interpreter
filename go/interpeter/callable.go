package interpeter

import "github.com/lukas-reining/lox/parser/statements"

type Callable interface {
	Stringifyable
	Arity() int
	Call(interpreter *Interpreter, args []LoxValue) (LoxValue, RuntimeError)
}

type LoxCallable struct {
	Callable
	arity int
	call  func(interpreter *Interpreter, args []LoxValue) (LoxValue, RuntimeError)
}

func (c *LoxCallable) ToString() string {
	return "<callable>"
}

func (c *LoxCallable) Arity() int {
	return c.arity
}

func (c *LoxCallable) Call(interpreter *Interpreter, args []LoxValue) (LoxValue, RuntimeError) {
	return c.call(interpreter, args)
}

func NewLoxCallable(arity int, call func(interpreter *Interpreter, args []LoxValue) (LoxValue, RuntimeError)) *LoxCallable {
	return &LoxCallable{
		arity: arity,
		call:  call,
	}
}

type LoxFunction struct {
	LoxCallable
	isInitializer bool
	closure       *Environment
	declaration   statements.Function[LoxValue, RuntimeError]
}

func NewLoxFunction(declaration statements.Function[LoxValue, RuntimeError], closure *Environment, isInitializer bool) *LoxFunction {
	arity := len(declaration.Params)

	return &LoxFunction{
		LoxCallable{
			arity: arity,
		},
		isInitializer,
		closure,
		declaration,
	}
}

func (f *LoxFunction) Call(interpreter *Interpreter, args []LoxValue) (LoxValue, RuntimeError) {
	environment := NewEnvironment(f.closure)

	for index, token := range f.declaration.Params {
		environment.define(token.Lexeme, args[index])
	}

	value, err := interpreter.executeBlock(f.declaration.Body, environment)
	if err != nil {
		return nil, err
	}

	if f.isInitializer {
		return f.closure.getAt(0, "this"), nil
	} else {
		return value, nil
	}
}

func (f *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	environment := NewEnvironment(f.closure)
	environment.define("this", instance)
	return NewLoxFunction(f.declaration, environment, f.isInitializer)
}

func (f *LoxFunction) ToString() string {
	return "<fn " + f.declaration.Name.Lexeme + ">"
}
