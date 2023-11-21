package parser

import "fmt"

type ParseError interface {
	error
	Error() string
	Line() int
	Lexeme() string
	Message() string
}

type BaseParseError struct {
	ParseError

	line    int
	message string
	lexeme  string
}

func (e *BaseParseError) Error() string {
	return fmt.Sprintf("[line %d] ParserError: %v", e.Line(), e.Message())
}

func (e *BaseParseError) Line() int {
	return e.line
}

func (e *BaseParseError) Lexeme() string {
	return e.lexeme
}

func (e *BaseParseError) Message() string {
	return e.message
}

func NewParseError(line int, lexeme string, message string) ParseError {
	return &BaseParseError{
		line: line, lexeme: lexeme, message: message,
	}
}
