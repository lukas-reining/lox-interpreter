package interpeter

import (
	"fmt"
	"github.com/lukas-reining/lox/parser/expressions"
	"github.com/lukas-reining/lox/parser/statements"
	"github.com/lukas-reining/lox/scanner"
)

type FunctionType string

const (
	NONE_FUNCTION FunctionType = "None"
	FUNCTION      FunctionType = "Function"
	METHOD        FunctionType = "Method"
	INITIALIZER   FunctionType = "Initializer"
)

type ClassType string

const (
	NONE_CLASS ClassType = "None"
	CLASS      ClassType = "Class"
)

type Resolver struct {
	e expressions.Visitor[LoxValue, RuntimeError]
	s statements.Visitor[LoxValue, RuntimeError]

	scopes              []map[string]bool
	currentFunctionType FunctionType
	currentClassType    ClassType
	interpreter         *Interpreter
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:         interpreter,
		currentClassType:    NONE_CLASS,
		currentFunctionType: NONE_FUNCTION,
	}
}

func (r *Resolver) VisitGroupingExpression(exp *expressions.Grouping[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	return nil, r.resolveExpression(exp.Exp)
}

func (r *Resolver) VisitBinaryExpression(exp *expressions.Binary[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	if err := r.resolveExpression(exp.Left); err != nil {
		return nil, err
	}

	if err := r.resolveExpression(exp.Right); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitUnaryExpression(exp *expressions.Unary[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	return nil, r.resolveExpression(exp.Right)
}

func (r *Resolver) VisitLiteralExpression(exp *expressions.Literal[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	return nil, nil
}

func (r *Resolver) VisitExpressionStatement(statement *statements.Expression[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	return nil, r.resolveExpression(statement.Exp)
}

func (r *Resolver) VisitPrintStatement(statement *statements.Print[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	return nil, r.resolveExpression(statement.Exp)
}

func (r *Resolver) VisitVarStatement(statement *statements.Var[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	if err := r.declare(statement.Name); err != nil {
		return nil, err
	}

	if statement.Initializer != nil {
		if err := r.resolveExpression(statement.Initializer); err != nil {
			return nil, err
		}
	}
	r.define(statement.Name)
	return nil, nil
}

func (r *Resolver) VisitVariableExpression(exp *expressions.Variable[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	if r.hasScopes() && r.declaredInScope(exp.Name) && r.currentScope()[exp.Name.Lexeme] == false {
		return nil, NewRuntimeError(exp.Name.Line, "Can't read local variable in its own initializer.")
	}

	r.resolveLocal(exp, exp.Name)
	return nil, nil
}

func (r *Resolver) VisitAssignmentExpression(exp *expressions.Assignment[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	if err := r.resolveExpression(exp.Value); err != nil {
		return nil, err
	}

	r.resolveLocal(exp, exp.Name)
	return nil, nil
}

func (r *Resolver) VisitBlockStatement(statement *statements.Block[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	r.beginScope()
	if err := r.Resolve(statement.Statements); err != nil {
		return nil, err
	}

	r.endScope()
	return nil, nil
}

func (r *Resolver) VisitIfStatement(statement *statements.If[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	if err := r.resolveExpression(statement.Condition); err != nil {
		return nil, err
	}

	if err := r.resolveStatement(statement.IfBranch); err != nil {
		return nil, err
	}

	if statement.ElseBranch == nil {
		return nil, nil
	}

	if err := r.resolveStatement(statement.IfBranch); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitLogicalExpression(expr *expressions.Logical[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	if err := r.resolveExpression(expr.Left); err != nil {
		return nil, err
	}

	if err := r.resolveExpression(expr.Right); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitWhileStatement(statement *statements.While[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	if err := r.resolveExpression(statement.Condition); err != nil {
		return nil, err
	}

	if err := r.resolveStatement(statement.Body); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitCallExpression(exp *expressions.Call[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	if err := r.resolveExpression(exp.Callee); err != nil {
		return nil, err
	}

	for _, arg := range exp.Params {
		if err := r.resolveExpression(arg); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitFunctionStatement(statement *statements.Function[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	if err := r.declare(statement.Name); err != nil {
		return nil, err
	}
	r.define(statement.Name)
	return nil, r.resolveFunction(*statement, FUNCTION)
}

func (r *Resolver) VisitClassStatement(statement *statements.Class[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	enclosingClassType := r.currentClassType
	r.currentClassType = CLASS

	if err := r.declare(statement.Name); err != nil {
		return nil, err
	}
	r.define(statement.Name)
	r.beginScope()

	r.currentScope()["this"] = true

	for _, method := range statement.Methods {
		declaration := METHOD
		if method.Name.Lexeme == "init" {
			declaration = INITIALIZER
		}

		if err := r.resolveFunction(method, declaration); err != nil {
			return nil, err
		}
	}

	r.endScope()
	r.currentClassType = enclosingClassType
	return nil, nil
}

func (r *Resolver) VisitGetExpression(exp *expressions.Get[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	return nil, r.resolveExpression(exp.Object)
}

func (r *Resolver) VisitSetExpression(exp *expressions.Set[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	if err := r.resolveExpression(exp.Object); err != nil {
		return nil, err
	}

	if err := r.resolveExpression(exp.Value); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitThisExpression(exp *expressions.This[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	if r.currentFunctionType == NONE_FUNCTION {
		return nil, NewRuntimeError(exp.Keyword.Line, "Can't use 'this' outside of a class.")
	}

	r.resolveLocal(exp, exp.Keyword)
	return nil, nil
}

func (r *Resolver) VisitReturnStatement(statement *statements.Return[LoxValue, RuntimeError]) (LoxValue, RuntimeError) {
	if r.currentFunctionType == NONE_FUNCTION {
		return nil, NewRuntimeError(statement.Keyword.Line, "Can't return from top-level code.")
	}

	if statement.Value != nil {
		if r.currentFunctionType == INITIALIZER {
			return nil, NewRuntimeError(statement.Keyword.Line, "Can't return a value from an initializer.")
		}

		return nil, r.resolveExpression(statement.Value)
	}

	return nil, nil
}

func (r *Resolver) resolveExpression(expression expressions.Expression[LoxValue, RuntimeError]) RuntimeError {
	_, err := expression.Accept(r)
	return err
}

func (r *Resolver) resolveFunction(statement statements.Function[LoxValue, RuntimeError], fnType FunctionType) RuntimeError {
	enclosingFunctionType := r.currentFunctionType
	r.currentFunctionType = fnType

	r.beginScope()

	for _, param := range statement.Params {
		if err := r.declare(param); err != nil {
			return err
		}

		r.define(param)
	}

	if err := r.Resolve(statement.Body); err != nil {
		return err
	}

	r.endScope()
	r.currentFunctionType = enclosingFunctionType
	return nil
}

func (r *Resolver) resolveStatement(statement statements.Statement[LoxValue, RuntimeError]) RuntimeError {
	_, err := statement.Accept(r)
	return err
}

func (r *Resolver) resolveLocal(expression expressions.Expression[LoxValue, RuntimeError], name scanner.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, hasValue := r.scopes[i][name.Lexeme]; hasValue {
			r.interpreter.resolve(expression, len(r.scopes)-i-1)
			return
		}
	}
}

func (r *Resolver) Resolve(statements []statements.Statement[LoxValue, RuntimeError]) RuntimeError {
	for _, statement := range statements {
		if err := r.resolveStatement(statement); err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) hasScopes() bool {
	return len(r.scopes) != 0
}

func (r *Resolver) currentScope() map[string]bool {
	return r.scopes[len(r.scopes)-1]
}

func (r *Resolver) declaredInScope(variable scanner.Token) bool {
	_, ok := r.scopes[len(r.scopes)-1][variable.Lexeme]
	return ok
}

func (r *Resolver) declare(name scanner.Token) RuntimeError {
	if !r.hasScopes() {
		return nil
	}

	if r.declaredInScope(name) {
		return NewRuntimeError(name.Line, fmt.Sprintf("Already a variable '%s' in this scope.", name.Lexeme))
	}

	r.currentScope()[name.Lexeme] = false
	return nil
}

func (r *Resolver) define(name scanner.Token) {
	if !r.hasScopes() {
		return
	}

	r.currentScope()[name.Lexeme] = true
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
}
