package interpeter

import (
	"fmt"
	"github.com/lukas-reining/lox/scanner"
)

type LoxClass struct {
	Stringifyable
	Callable
	Name    string
	Methods map[string]LoxFunction
}

func NewLoxClass(name string, methods map[string]LoxFunction) *LoxClass {
	return &LoxClass{
		Name:    name,
		Methods: methods,
	}
}

func (c *LoxClass) Call(interpreter *Interpreter, args []LoxValue) (LoxValue, RuntimeError) {
	instance := NewLoxInstance(c)

	initializer, hasInit := c.Methods["init"]
	if hasInit {
		if _, err := initializer.Bind(instance).Call(interpreter, args); err != nil {
			return nil, err
		}
	}

	return instance, nil
}

func (c *LoxClass) Arity() int {
	initializer, hasInit := c.Methods["init"]
	if hasInit {
		return initializer.arity
	} else {
		return 0
	}
}

func (c *LoxClass) ToString() string {
	return c.Name
}

type LoxInstance struct {
	Stringifyable
	Class  *LoxClass
	Fields map[string]LoxValue
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		Class:  class,
		Fields: make(map[string]LoxValue),
	}
}

func (c *LoxInstance) get(name scanner.Token) (LoxValue, RuntimeError) {
	if value, hasValue := c.Fields[name.Lexeme]; hasValue {
		return value, nil
	}

	if value, hasValue := c.Class.Methods[name.Lexeme]; hasValue {
		scopedMethod := value.Bind(c)
		return scopedMethod, nil
	}

	return NewRuntimeError(name.Line, fmt.Sprintf("Undefined property '%s'.", name.Lexeme)), nil
}

func (c *LoxInstance) set(name scanner.Token, value LoxValue) {
	c.Fields[name.Lexeme] = value
}

func (c *LoxInstance) ToString() string {
	return c.Class.Name + " instance"
}
