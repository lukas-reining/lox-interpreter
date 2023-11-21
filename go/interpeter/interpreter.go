package interpeter

import (
	"fmt"
	"github.com/lukas-reining/lox/parser/expressions"
	"github.com/lukas-reining/lox/parser/statements"
	"github.com/lukas-reining/lox/scanner"
	"time"
)

type Interpreter struct {
	e expressions.Visitor[LoxValue, RuntimeError]
	s statements.Visitor[LoxValue, RuntimeError]

	env     *Environment
	globals *Environment
	locals  map[expressions.Expression[LoxValue, RuntimeError]]int
}

func NewInterpreterWithEnv(env *Environment) Interpreter {
	newEnv := NewEnvironment(env) // TODO Make work for REPL
	defineGlobals(newEnv)

	return Interpreter{
		globals: newEnv,
		env:     newEnv,
		locals:  map[expressions.Expression[LoxValue, RuntimeError]]int{},
	}
}

func NewInterpreter() Interpreter {
	return NewInterpreterWithEnv(nil)
}

func defineGlobals(globals *Environment) {
	globals.define("env", "LOX")

	globals.define("clock", NewLoxCallable(0, func(interpreter *Interpreter, args []LoxValue) (LoxValue, RuntimeError) {
		return float64(time.Now().UnixMicro()), nil
	}))

	globals.define("printStackDepth", NewLoxCallable(0, func(interpreter *Interpreter, args []LoxValue) (LoxValue, RuntimeError) {
		fmt.Printf("Stack Depth: %d\n", interpreter.env.level)
		return nil, nil
	}))
}

func GetGlobalEnv() *Environment {
	env := NewEnvironment(nil)
	defineGlobals(env)
	return env
}

func isTruthy(value LoxValue) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	default:
		return true
	}
}

func Stringify(value LoxValue) string {
	if value == nil {
		return fmt.Sprint("nil")
	} else {
		switch v := value.(type) {
		case string:
			return fmt.Sprint(value)
		case float64:
			return fmt.Sprint(value)
		case bool:
			return fmt.Sprint(value)
		case Stringifyable:
			return v.ToString()
		default:
			return "<Unknown>"
		}
	}
}

func createBinaryNumberOperatorError(operator scanner.Token) RuntimeError {
	return NewRuntimeError(operator.Line, "Operands must be numbers.")
}

func createBinaryNumberOrStringOperatorError(operator scanner.Token) RuntimeError {
	return NewRuntimeError(operator.Line, "Operands must be two numbers or two strings.")
}

func (i *Interpreter) VisitGroupingExpression(exp *expressions.Grouping[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	return i.evaluate(exp.Exp)
}

func (i *Interpreter) VisitBinaryExpression(exp *expressions.Binary[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	left, err := i.evaluate(exp.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(exp.Right)
	if err != nil {
		return nil, err
	}

	leftNumber, leftNumberOk := left.(float64)
	rightNumber, rightNumberOk := right.(float64)

	leftString, leftStringOk := left.(string)
	rightString, rightStringOk := right.(string)

	switch exp.Operator.Type {
	case scanner.MINUS:
		if leftNumberOk && rightNumberOk {
			return leftNumber - rightNumber, nil
		}
		return nil, createBinaryNumberOperatorError(exp.Operator)
	case scanner.PLUS:
		if leftNumberOk && rightNumberOk {
			return leftNumber + rightNumber, nil
		}
		if leftStringOk && rightStringOk {
			return leftString + rightString, nil
		}
		if leftNumberOk && rightStringOk {
			return Stringify(leftNumber) + rightString, nil
		}
		if leftStringOk && rightNumberOk {
			return leftString + Stringify(rightNumber), nil
		}
		return nil, createBinaryNumberOrStringOperatorError(exp.Operator)
	case scanner.STAR:
		if leftNumberOk && rightNumberOk {
			return leftNumber * rightNumber, nil
		}
		return nil, createBinaryNumberOperatorError(exp.Operator)
	case scanner.SLASH:
		if leftNumberOk && rightNumberOk {
			return leftNumber / rightNumber, nil
		}
		return nil, createBinaryNumberOperatorError(exp.Operator)
	case scanner.GREATER:
		if leftNumberOk && rightNumberOk {
			return leftNumber > rightNumber, nil
		}
		return nil, createBinaryNumberOperatorError(exp.Operator)
	case scanner.GREATER_EQUAL:
		if leftNumberOk && rightNumberOk {
			return leftNumber >= rightNumber, nil
		}
		return nil, createBinaryNumberOperatorError(exp.Operator)
	case scanner.LESS:
		if leftNumberOk && rightNumberOk {
			return leftNumber < rightNumber, nil
		}
		return nil, createBinaryNumberOperatorError(exp.Operator)
	case scanner.LESS_EQUAL:
		if leftNumberOk && rightNumberOk {
			return leftNumber <= rightNumber, nil
		}
		return nil, createBinaryNumberOperatorError(exp.Operator)
	case scanner.EQUAL_EQUAL:
		return left == right, nil
	case scanner.BANG_EQUAL:
		return left != right, nil
	}

	// Unreachable.
	return nil, nil
}

func (i *Interpreter) VisitUnaryExpression(exp *expressions.Unary[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	right, err := i.evaluate(exp.Right)

	if err != nil {
		return nil, err
	}

	switch exp.Operator.Type {
	case scanner.MINUS:
		if number, ok := right.(float64); ok {
			return -number, nil
		}
	case scanner.BANG:
		return isTruthy(right), nil
	}

	// Unreachable.
	return nil, nil
}

func (i *Interpreter) VisitLiteralExpression(exp *expressions.Literal[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	return exp.Literal, nil
}

func (i *Interpreter) VisitExpressionStatement(statement *statements.Expression[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	return i.evaluate(statement.Exp)
}

func (i *Interpreter) VisitPrintStatement(statement *statements.Print[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	value, err := i.evaluate(statement.Exp)

	if err == nil {
		println(Stringify(value))
	}

	return nil, err
}

func (i *Interpreter) VisitVarStatement(statement *statements.Var[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	if statement.Initializer != nil {
		value, err := i.evaluate(statement.Initializer)

		if err != nil {
			return nil, err
		}

		i.env.define(statement.Name.Lexeme, value)
	} else {
		i.env.define(statement.Name.Lexeme, nil)
	}
	return nil, nil
}

func (i *Interpreter) VisitVariableExpression(exp *expressions.Variable[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	return i.lookupVariable(exp.Name, exp)
}

func (i *Interpreter) VisitAssignmentExpression(exp *expressions.Assignment[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	value, err := i.evaluate(exp.Value)
	if err != nil {
		return nil, err
	}

	distance, hasDistance := i.locals[exp]

	if hasDistance {
		i.env.assignAt(distance, exp.Name, value)
	} else {
		if err := i.globals.assign(exp.Name, value); err != nil {
			return nil, err
		}
	}

	return value, nil
}

func (i *Interpreter) VisitBlockStatement(statement *statements.Block[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	currentEnv := i.env
	newEnv := NewEnvironment(currentEnv)
	return i.executeBlock(statement.Statements, newEnv)
}

func (i *Interpreter) VisitIfStatement(statement *statements.If[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	condition, err := i.evaluate(statement.Condition)

	if err != nil {
		return nil, err
	}

	var value LoxValue
	var evalErr RuntimeError
	if isTruthy(condition) && statement.IfBranch != nil {
		value, evalErr = i.execute(statement.IfBranch)
	} else if statement.ElseBranch != nil {
		value, evalErr = i.execute(statement.ElseBranch)
	}
	return value, evalErr
}

func (i *Interpreter) VisitLogicalExpression(expr *expressions.Logical[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	left, err := i.evaluate(expr.Left)

	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == scanner.OR && isTruthy(left) {
		return left, nil
	} else if expr.Operator.Type != scanner.OR && !isTruthy(left) {
		return left, nil
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitWhileStatement(statement *statements.While[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	condition, err := i.evaluate(statement.Condition)

	for err == nil && isTruthy(condition) {
		_, bodyError := i.execute(statement.Body)

		if bodyError != nil {
			return nil, bodyError
		}

		condition, err = i.evaluate(statement.Condition)
	}

	return nil, err
}

func (i *Interpreter) VisitCallExpression(exp *expressions.Call[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	callee, err := i.evaluate(exp.Callee)
	if err != nil {
		return nil, err
	}

	args := []LoxValue{}
	for _, arg := range exp.Params {
		if argValue, err := i.evaluate(arg); err != nil {
			return nil, err
		} else {
			args = append(args, argValue)
		}
	}

	callable, isCallable := callee.(Callable)
	if !isCallable {
		return nil, NewRuntimeError(exp.Parenthesis.Line, "Can only call functions and classes.")
	}

	if len(args) != callable.Arity() {
		return nil, NewRuntimeError(exp.Parenthesis.Line, fmt.Sprintf("Expected %d arguments but got %d.", callable.Arity(), len(args)))
	}

	return callable.Call(i, args)
}

func (i *Interpreter) VisitGetExpression(exp *expressions.Get[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	object, err := i.evaluate(exp.Object)
	if err != nil {
		return nil, err
	}

	switch o := object.(type) {
	case *LoxInstance:
		value, err := o.get(exp.Name)
		return value, err
	case LoxInstance:
		value, err := o.get(exp.Name)
		return value, err
	}

	return nil, NewRuntimeError(exp.Name.Line, "Only instances have properties.")
}

func (i *Interpreter) VisitSetExpression(exp *expressions.Set[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	object, err := i.evaluate(exp.Object)
	if err != nil {
		return nil, err
	}

	value, err := i.evaluate(exp.Value)
	if err != nil {
		return nil, err
	}

	switch o := object.(type) {
	case *LoxInstance:
		o.set(exp.Name, value)
		return nil, nil
	case LoxInstance:
		o.set(exp.Name, value)
		return nil, nil
	}

	return nil, NewRuntimeError(exp.Name.Line, "Only instances have Fields.")
}

func (i *Interpreter) VisitFunctionStatement(statement *statements.Function[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	function := NewLoxFunction(*statement, i.env, false)
	i.env.define(function.declaration.Name.Lexeme, function)
	return nil, nil
}

func (i *Interpreter) VisitClassStatement(statement *statements.Class[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	i.env.define(statement.Name.Lexeme, nil)

	methods := map[string]LoxFunction{}
	for _, method := range statement.Methods {
		fn := *NewLoxFunction(method, i.env, method.Name.Lexeme == "init")
		methods[method.Name.Lexeme] = fn
	}

	class := NewLoxClass(statement.Name.Lexeme, methods)
	return nil, i.env.assign(statement.Name, class)
}

func (i *Interpreter) VisitReturnStatement(statement *statements.Return[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	if statement.Value == nil {
		return NewReturnValue(nil), nil
	}

	if value, err := i.evaluate(statement.Value); err != nil {
		return nil, err
	} else {
		return NewReturnValue(value), nil
	}
}

func (i *Interpreter) VisitThisExpression(exp *expressions.This[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	return i.lookupVariable(exp.Keyword, exp)
}

func (i *Interpreter) lookupVariable(name scanner.Token, exp expressions.Expression[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	distance, hasDistance := i.locals[exp]

	if hasDistance {
		return i.env.getAt(distance, name.Lexeme), nil
	} else {
		return i.globals.get(name)
	}
}

func (i *Interpreter) executeBlock(statements []statements.Statement[LoxValue, RuntimeError], env *Environment) (LoxValue, RuntimeError) {
	previousEnv := i.env
	i.env = env

	for _, statement := range statements {
		value, err := i.execute(statement)

		if err != nil {
			i.env = previousEnv
			return nil, err
		}

		switch returnValue := value.(type) {
		case ReturnValue:
			i.env = previousEnv
			return returnValue.Value, nil
		case *ReturnValue:
			i.env = previousEnv
			return returnValue.Value, nil
		}
	}

	i.env = previousEnv
	return nil, nil
}

func (i *Interpreter) evaluate(exp expressions.Expression[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	return exp.Accept(i)
}

func (i *Interpreter) resolve(exp expressions.Expression[LoxValue, RuntimeError], depth int) {
	i.locals[exp] = depth
}

func (i *Interpreter) execute(statement statements.Statement[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	return statement.Accept(i)
}

func (i *Interpreter) Interpret(statements []statements.Statement[LoxValue, RuntimeError]) (LoxValue, *Environment, RuntimeError) {
	var lastValue LoxValue

	for _, statement := range statements {
		value, err := i.execute(statement)
		lastValue = value

		if err != nil {
			return nil, i.env, err
		}
	}

	return lastValue, i.env, nil
}
