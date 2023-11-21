package parser

import (
	"fmt"
	"github.com/lukas-reining/lox/interpeter"
	"github.com/lukas-reining/lox/parser/expressions"
	"github.com/lukas-reining/lox/parser/statements"
	"github.com/lukas-reining/lox/scanner"
)

type Parser[T any, Err error] struct {
	Tokens  []scanner.Token
	current int
}

func NewParser[T any, Err error](tokens []scanner.Token) *Parser[T, Err] {
	return &Parser[T, Err]{Tokens: tokens}
}

func (p *Parser[T, Err]) isAtEnd() bool {
	return p.peek().Type == scanner.EOF
}

func (p *Parser[T, Err]) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *Parser[T, Err]) previous() scanner.Token {
	return p.Tokens[p.current-1]
}

func (p *Parser[T, Err]) peek() scanner.Token {
	return p.Tokens[p.current]
}

func (p *Parser[T, Err]) check(tokenType scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokenType
}

func (p *Parser[T, Err]) match(types ...scanner.TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser[T, Err]) consume(tokenType scanner.TokenType, message string) (scanner.Token, ParseError) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return scanner.Token{}, NewParseError(p.previous().Line, "", message)
}

func (p *Parser[T, Err]) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == scanner.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case scanner.CLASS:
		case scanner.FUN:
		case scanner.VAR:
		case scanner.FOR:
		case scanner.IF:
		case scanner.WHILE:
		case scanner.PRINT:
		case scanner.RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser[T, Err]) primary() (expressions.Expression[T, Err], ParseError) {
	if p.match(scanner.FALSE) {
		return expressions.NewLiteral[T, Err](false), nil
	}

	if p.match(scanner.TRUE) {
		return expressions.NewLiteral[T, Err](true), nil
	}

	if p.match(scanner.NIL) {
		return expressions.NewLiteral[T, Err](nil), nil
	}

	if p.match(scanner.NUMBER, scanner.STRING) {
		return expressions.NewLiteral[T, Err](p.previous().Literal), nil
	}

	if p.match(scanner.LEFT_PAREN) {
		if expr, err := p.expression(); err != nil {
			return nil, err
		} else if _, err := p.consume(scanner.RIGHT_PAREN, "Expect ')' after expression."); err == nil {
			return expressions.NewGrouping[T, Err](expr), nil
		} else {
			return nil, err
		}
	}

	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.previous()

		if right, err := p.unary(); err != nil {
			return nil, err
		} else if _, err := p.consume(scanner.RIGHT_PAREN, "Expect ')' after expression."); err == nil {
			return expressions.NewUnary[T, Err](operator, right), nil
		} else {
			return nil, err
		}
	}

	if p.match(scanner.THIS) {
		return expressions.NewThis[T, Err](p.previous()), nil
	}

	if p.match(scanner.IDENTIFIER) {
		return expressions.NewVariable[T, Err](p.previous()), nil
	}

	return nil, NewParseError(p.current, p.peek().Lexeme, "Expected expression!")
}

func (p *Parser[T, Err]) unary() (expressions.Expression[T, Err], ParseError) {

	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.previous()

		if right, err := p.unary(); err != nil {
			return nil, err
		} else {
			return expressions.NewUnary(operator, right), nil
		}
	}

	return p.call()
}

func (p *Parser[T, Err]) factor() (expressions.Expression[T, Err], ParseError) {
	left, err := p.unary()

	if err != nil {
		return nil, err
	}

	for p.match(scanner.SLASH, scanner.STAR) {
		operator := p.previous()
		if right, err := p.unary(); err != nil {
			return nil, err
		} else {
			left = expressions.NewBinary(left, operator, right)
		}

	}

	return left, nil
}

func (p *Parser[T, Err]) term() (expressions.Expression[T, Err], ParseError) {
	left, err := p.factor()

	if err != nil {
		return nil, err
	}

	for p.match(scanner.MINUS, scanner.PLUS) {
		operator := p.previous()

		if right, err := p.factor(); err != nil {
			return nil, err
		} else {
			left = expressions.NewBinary(left, operator, right)
		}
	}

	return left, nil
}

func (p *Parser[T, Err]) comparison() (expressions.Expression[T, Err], ParseError) {
	left, err := p.term()

	if err != nil {
		return nil, err
	}

	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.previous()

		if right, err := p.term(); err != nil {
			return nil, err
		} else {
			left = expressions.NewBinary(left, operator, right)
		}
	}

	return left, nil
}

func (p *Parser[T, Err]) equality() (expressions.Expression[T, Err], ParseError) {
	left, err := p.comparison()

	if err != nil {
		return nil, err
	}

	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		operator := p.previous()

		if right, err := p.comparison(); err != nil {
			return nil, err
		} else {
			left = expressions.NewBinary(left, operator, right)
		}
	}

	return left, nil
}

func (p *Parser[T, Err]) and() (expressions.Expression[T, Err], ParseError) {
	expr, err := p.equality()

	if err != nil {
		return nil, err
	}

	for p.match(scanner.AND) {
		operator := p.previous()

		right, err := p.equality()

		if err != nil {
			return nil, err
		}

		return expressions.NewLogical(expr, operator, right), nil
	}

	return expr, nil
}

func (p *Parser[T, Err]) or() (expressions.Expression[T, Err], ParseError) {
	expr, err := p.and()

	if err != nil {
		return nil, err
	}

	for p.match(scanner.OR) {
		operator := p.previous()

		right, err := p.and()

		if err != nil {
			return nil, err
		}

		return expressions.NewLogical(expr, operator, right), nil
	}

	return expr, nil
}

func (p *Parser[T, Err]) assignment() (expressions.Expression[T, Err], ParseError) {
	exp, err := p.or()

	if err != nil {
		return nil, err
	}

	if p.match(scanner.EQUAL) {
		token := p.previous()
		value, err := p.assignment()

		if err != nil {
			return nil, err
		}

		switch expression := exp.(type) {
		case *expressions.Variable[T, Err]:
			name := expression.Name
			return expressions.NewAssignment(name, value), nil
		case *expressions.Get[T, Err]:
			return expressions.NewSet(expression.Object, expression.Name, value), nil
		}

		return nil, NewParseError(token.Line, token.Lexeme, "Invalid assignment target.")
	}

	return exp, nil
}

func (p *Parser[T, Err]) expression() (expressions.Expression[T, Err], ParseError) {
	return p.assignment()
}

func (p *Parser[T, Err]) expressionStatement() (statements.Statement[T, Err], ParseError) {
	value, err := p.expression()

	if err != nil {
		return nil, err
	}

	if _, err := p.consume(scanner.SEMICOLON, "Expect ';' after value."); err == nil {
		return statements.NewExpression(value), nil
	} else {
		return nil, err
	}
}

func (p *Parser[T, Err]) printStatement() (statements.Statement[T, Err], ParseError) {
	value, err := p.expression()

	if err != nil {
		return nil, err
	}

	if _, err := p.consume(scanner.SEMICOLON, "Expect ';' after value."); err == nil {
		return statements.NewPrintStatement(value), nil
	} else {
		return nil, err
	}
}

func (p *Parser[T, Err]) block() ([]statements.Statement[T, Err], ParseError) {
	var stmnts []statements.Statement[T, Err]

	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		if declaration, err := p.declaration(); err != nil {
			return nil, err
		} else {
			stmnts = append(stmnts, declaration)
		}
	}

	if _, err := p.consume(scanner.RIGHT_BRACE, "Expect '}' after block."); err != nil {
		return nil, err
	}

	return stmnts, nil
}

func (p *Parser[T, Err]) blockStatement() (statements.Statement[T, Err], ParseError) {
	if block, err := p.block(); err != nil {
		return nil, err
	} else {
		return statements.NewBlock(block), nil
	}

}

func (p *Parser[T, Err]) ifStatement() (statements.Statement[T, Err], ParseError) {
	if _, err := p.consume(scanner.LEFT_PAREN, "Expect '(' after 'if'."); err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err = p.consume(scanner.RIGHT_PAREN, "Expect ')' after if condition."); err != nil {
		return nil, err
	}

	ifBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch statements.Statement[T, Err]
	if p.match(scanner.ELSE) {
		statement, err := p.statement()

		if err != nil {
			return nil, err
		}

		elseBranch = statement
	}

	return statements.NewIf(condition, ifBranch, elseBranch), nil
}

func (p *Parser[T, Err]) whileStatement() (statements.Statement[T, Err], ParseError) {

	if _, err := p.consume(scanner.LEFT_PAREN, "Expect '(' after 'while'."); err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err = p.consume(scanner.RIGHT_PAREN, "Expect ')' after condition."); err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return statements.NewWhile(condition, body), nil
}

func (p *Parser[T, Err]) forStatement() (statements.Statement[T, Err], ParseError) {
	if _, err := p.consume(scanner.LEFT_PAREN, "Expect '(' after 'for'."); err != nil {
		return nil, err
	}

	var initializer statements.Statement[T, Err]
	var err ParseError
	if p.match(scanner.SEMICOLON) {
		initializer = nil
	} else if p.match(scanner.VAR) {
		initializer, err = p.varDeclaration()
	} else {
		initializer, err = p.expressionStatement()
	}

	if err != nil {
		return nil, err
	}

	var condition expressions.Expression[T, Err]
	if !p.check(scanner.SEMICOLON) {
		condition, err = p.expression()
	}

	if err != nil {
		return nil, err
	}

	if _, err = p.consume(scanner.SEMICOLON, "Expect ';' after loop condition."); err != nil {
		return nil, err
	}

	var increment expressions.Expression[T, Err]
	if !p.check(scanner.RIGHT_PAREN) {
		increment, err = p.expression()
	}

	if _, err = p.consume(scanner.RIGHT_PAREN, "Expect ')' after for clauses."); err != nil {
		return nil, err
	}

	body, err := p.statement()

	if increment != nil {
		blockStatements := []statements.Statement[T, Err]{body, statements.NewExpression(increment)}
		body = statements.NewBlock(blockStatements)
	}

	if condition == nil {
		condition = expressions.NewLiteral[T, Err](true)
	}

	body = statements.NewWhile(condition, body)

	if initializer != nil {
		blockStatements := []statements.Statement[T, Err]{initializer, body}
		body = statements.NewBlock(blockStatements)
	}

	return body, nil
}

func (p *Parser[T, Err]) returnStatement() (statements.Statement[T, Err], ParseError) {
	keyword := p.previous()

	var value expressions.Expression[T, Err]
	var err ParseError

	if !p.check(scanner.SEMICOLON) {
		value, err = p.expression()
	}

	if err != nil {
		return nil, err
	}

	if _, err := p.consume(scanner.SEMICOLON, "Expect ';' after return value."); err != nil {
		return nil, err
	}

	return statements.NewReturn(keyword, value), nil
}

func (p *Parser[T, Err]) statement() (statements.Statement[T, Err], ParseError) {
	if p.match(scanner.FOR) {
		return p.forStatement()
	}

	if p.match(scanner.IF) {
		return p.ifStatement()
	}

	if p.match(scanner.PRINT) {
		return p.printStatement()
	}

	if p.match(scanner.RETURN) {
		return p.returnStatement()
	}

	if p.match(scanner.WHILE) {
		return p.whileStatement()
	}

	if p.match(scanner.LEFT_BRACE) {
		return p.blockStatement()
	}

	return p.expressionStatement()

}

func (p *Parser[T, Err]) varDeclaration() (statements.Statement[T, Err], ParseError) {
	name, consumeEr := p.consume(scanner.IDENTIFIER, "Expect variable name.")

	if consumeEr != nil {
		return nil, consumeEr
	}

	var initializer expressions.Expression[T, Err]
	var err ParseError
	if p.match(scanner.EQUAL) {
		initializer, err = p.expression()
	}

	if err != nil {
		return nil, err
	}

	if _, err := p.consume(scanner.SEMICOLON, "Expect ';' after variable declaration."); err == nil {
		return statements.NewVar(name, initializer), nil
	} else {
		return nil, err
	}
}

func (p *Parser[T, Err]) function(kind interpeter.FunctionType) (*statements.Function[T, Err], ParseError) {
	name, err := p.consume(scanner.IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(scanner.LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name.", kind)); err != nil {
		return nil, err
	}

	var params []scanner.Token
	hasParamArg := !p.check(scanner.RIGHT_PAREN)
	for hasParamArg {
		if len(params) >= 255 {
			return nil, NewParseError(p.peek().Line, "", "Can't have more than 255 parameters.")
		}

		if param, err := p.consume(scanner.IDENTIFIER, "Expect parameter name."); err != nil {
			return nil, err
		} else {
			params = append(params, param)
		}

		hasParamArg = p.match(scanner.COMMA)
	}

	if _, err := p.consume(scanner.RIGHT_PAREN, fmt.Sprintf("Expect ')' after %s parameters.", kind)); err != nil {
		return nil, err
	}

	if _, err := p.consume(scanner.LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body.", kind)); err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return statements.NewFunction(name, params, body), nil
}

func (p *Parser[T, Err]) classDeclaration() (statements.Statement[T, Err], ParseError) {
	name, err := p.consume(scanner.IDENTIFIER, "Expect class name.")
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(scanner.LEFT_BRACE, "Expect '{' before class body."); err != nil {
		return nil, err
	}

	var methods []statements.Function[T, Err]
	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		if declaration, err := p.function(interpeter.METHOD); err != nil {
			return nil, err
		} else {
			methods = append(methods, *declaration)
		}
	}

	if _, err := p.consume(scanner.RIGHT_BRACE, "Expect '}' after class body."); err != nil {
		return nil, err
	}

	return statements.NewClass(name, methods), nil
}

func (p *Parser[T, Err]) declaration() (statements.Statement[T, Err], ParseError) {
	var value statements.Statement[T, Err]
	var err ParseError

	if p.match(scanner.CLASS) {
		return p.classDeclaration()
	} else if p.match(scanner.VAR) {
		value, err = p.varDeclaration()
	} else if p.match(scanner.FUN) {
		return p.function("function")
	} else if p.match(scanner.VAR) {
		value, err = p.varDeclaration()
	} else {
		value, err = p.statement()
	}

	if err != nil {
		p.synchronize()
		return nil, err
	}

	return value, nil
}

func (p *Parser[T, Err]) finishCall(expr expressions.Expression[T, Err]) (expressions.Expression[T, Err], ParseError) {
	var args []expressions.Expression[T, Err]

	hasNextArg := !p.check(scanner.RIGHT_PAREN)
	for hasNextArg {
		nextArg, err := p.expression()

		if err != nil {
			return nil, err
		}

		if len(args) >= 255 {
			return nil, NewParseError(p.peek().Line, "", "Can't have more than 255 arguments.")
		}

		args = append(args, nextArg)
		hasNextArg = p.match(scanner.COMMA)
	}

	if paren, err := p.consume(scanner.RIGHT_PAREN, "Expect ')' after arguments."); err != nil {
		return nil, err
	} else {
		return expressions.NewCall(expr, paren, args), nil
	}
}

func (p *Parser[T, Err]) call() (expressions.Expression[T, Err], ParseError) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(scanner.LEFT_PAREN) {
			expr, err = p.finishCall(expr)

			if err != nil {
				return nil, err
			}
		} else if p.match(scanner.DOT) {
			name, err := p.consume(scanner.IDENTIFIER, "Expect class name.")

			if err != nil {
				return nil, err
			}

			expr = expressions.NewGet(expr, name)
		} else {
			break
		}
	}

	return expr, err
}

func (p *Parser[T, Err]) Parse() ([]statements.Statement[T, Err], ParseError) {
	var stmnts []statements.Statement[T, Err]

	for !p.isAtEnd() {
		statement, err := p.declaration()

		if err != nil {
			return nil, err
		}

		stmnts = append(stmnts, statement)
	}

	return stmnts, nil
}
