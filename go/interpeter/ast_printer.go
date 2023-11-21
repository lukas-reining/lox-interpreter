package interpeter

import (
	"fmt"
	"github.com/lukas-reining/lox/parser/expressions"
	"strconv"
)

type AstPrinter struct {
	expressions.Visitor[string, RuntimeError]
}

func NewAstPrinter() AstPrinter {
	return AstPrinter{}
}

func (a *AstPrinter) Print(exp expressions.Expression[string, RuntimeError]) (string, RuntimeError) {
	return exp.Accept(a)
}

func (a *AstPrinter) VisitGroupingExpression(exp *expressions.Grouping[string, RuntimeError]) (string, RuntimeError) {
	return a.parenthesize("group", exp.Exp), nil
}

func (a *AstPrinter) VisitBinaryExpression(exp *expressions.Binary[string, RuntimeError]) (string, RuntimeError) {
	return a.parenthesize(exp.Operator.Lexeme, exp.Left, exp.Right), nil
}

func (a *AstPrinter) VisitUnaryExpression(exp *expressions.Unary[string, RuntimeError]) (string, RuntimeError) {
	return a.parenthesize(exp.Operator.Lexeme, exp.Right), nil
}

func (a *AstPrinter) VisitLiteralExpression(exp *expressions.Literal[string, RuntimeError]) (string, RuntimeError) {
	switch value := exp.Literal.(type) {
	case string:
		return value, nil
	case int:
		return strconv.Itoa(value), nil
	case float32:
		return fmt.Sprintf("%f", value), nil
	case float64:
		return fmt.Sprintf("%f", value), nil
	default:
		return "nil", nil
	}
}

func (a *AstPrinter) parenthesize(name string, expressions ...expressions.Expression[string, RuntimeError]) string {
	result := "(" + name

	for _, exp := range expressions {
		if value, err := exp.Accept(a); err != nil {
			result += " "
			result += value
		}
	}

	result += ")"
	return result
}
