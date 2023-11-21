package interpeter

import (
	"github.com/lukas-reining/lox/scanner"
)

type Environment struct {
	values    map[string]LoxValue
	level     int
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	level := 0
	if enclosing != nil {
		level = enclosing.level + 1
	}

	return &Environment{
		values:    make(map[string]LoxValue),
		level:     level,
		enclosing: enclosing,
	}
}

func (e *Environment) hasParent() bool {
	return e.enclosing != nil
}

func (e *Environment) define(name string, value LoxValue) {
	e.values[name] = value
}

func (e *Environment) assign(name scanner.Token, value LoxValue) RuntimeError {
	if !e.exists(name) {
		if e.enclosing != nil {
			return e.enclosing.assign(name, value)
		}

		return NewRuntimeError(name.Line, "Undefined variable '"+name.Lexeme+"'.")
	}

	e.values[name.Lexeme] = value
	return nil
}

func (e *Environment) assignAt(distance int, name scanner.Token, value LoxValue) {
	env := e.ancestor(distance)
	env.values[name.Lexeme] = value
}

func (e *Environment) exists(name scanner.Token) bool {
	_, ok := e.values[name.Lexeme]
	return ok
}

func (e *Environment) get(name scanner.Token) (LoxValue, RuntimeError) {
	value, ok := e.values[name.Lexeme]

	if ok {
		return value, nil
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	return nil, NewRuntimeError(name.Line, "Undefined variable '"+name.Lexeme+"'.")
}

func (e *Environment) getAt(distance int, name string) LoxValue {
	env := e.ancestor(distance)
	return env.values[name] // TODO Do we have to check here?
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e

	for i := 0; i < distance; i++ {
		env = env.enclosing
	}

	return env
}
