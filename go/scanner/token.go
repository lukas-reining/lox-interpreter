package scanner

type TokenType string

const (
	// Single-character tokens
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN           = "RIGHT_PAREN"
	LEFT_BRACE            = "LEFT_BRACE"
	RIGHT_BRACE           = "RIGHT_BRACE"
	COMMA                 = "COMMA"
	DOT                   = "DOT"
	MINUS                 = "MINUS"
	PLUS                  = "PLUS"
	SEMICOLON             = "SEMICOLON"
	SLASH                 = "SLASH"
	STAR                  = "STAR"

	// One or two character tokens
	BANG          = "BANG"
	BANG_EQUAL    = "BANG_EQUAL"
	EQUAL         = "EQUAL"
	EQUAL_EQUAL   = "EQUAL_EQUAL"
	GREATER       = "GREATER"
	GREATER_EQUAL = "GREATER_EQUAL"
	LESS          = "LESS"
	LESS_EQUAL    = "LESS_EQUAL"

	// Literals
	IDENTIFIER = "IDENTIFIER"
	STRING     = "STRING"
	NUMBER     = "NUMBER"

	// Keywords
	AND   = "AND"
	CLASS = "CLASS"
	ELSE  = "ELSE"
	FALSE = "FALSE"
	FUN   = "FUN"
	FOR   = "FOR"
	IF    = "IF"
	NIL   = "NIL"
	OR    = "OR"

	PRINT  = "PRINT"
	RETURN = "RETURN"
	SUPER  = "SUPER"
	THIS   = "THIS"
	TRUE   = "TRUE"
	VAR    = "VAR"
	WHILE  = "WHILE"
	EOF    = "EOF"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func NewToken(tokenType TokenType, lexeme string, literal any, line int) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t *Token) ToString() string {
	switch literal := t.Literal.(type) {
	case string:
		return string(t.Type) + " " + t.Lexeme + " " + literal
	default:
		return string(t.Type) + " " + t.Lexeme
	}
}
