package scanner

import (
	"strconv"
)

type Scanner struct {
	source  string
	start   int
	current int
	line    int
	tokens  []Token
}

func NewScanner(source string) Scanner {
	return Scanner{source: source, start: 0, current: 0, line: 1}
}

func (s *Scanner) Source() (source string) {
	return s.source
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() byte {
	current := s.source[s.current]
	s.current += 1
	return current
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}

	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}

	return s.source[s.current+1]
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != expected {
		return false
	} else {
		s.current += 1
		return true
	}
}

func (s *Scanner) matchOrElse(expected byte, onMatch TokenType, onNoMatch TokenType) TokenType {
	if s.match(expected) {
		return onMatch
	} else {
		return onNoMatch
	}
}

func (s *Scanner) isDigit(character byte) bool {
	return character >= '0' && character <= '9'
}

func (s *Scanner) isAlpha(character byte) bool {
	return (character >= 'a' && character <= 'z') ||
		(character >= 'A' && character <= 'Z') ||
		character == '_'
}

func (s *Scanner) isAlphaNumeric(character byte) bool {
	return s.isAlpha(character) || s.isDigit(character)
}

func (s *Scanner) handleString() ScannerError {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line += 1
		}
		s.advance()
	}

	if s.isAtEnd() {
		return NewScannerError(s.current, "Unterminated string.")
	}

	// The closing.
	s.advance()

	// Trim the surrounding quotes.
	value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, value)

	return nil
}

func (s *Scanner) handleNumber() ScannerError {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	if number, err := strconv.ParseFloat(s.source[s.start:s.current], 64); err == nil {
		s.addToken(NUMBER, number)
		return nil
	} else {
		return NewScannerError(s.current, "Could not parse number.")
	}
}

func (s *Scanner) handleIdentifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	var tokenType TokenType
	switch text {
	case "and":
		tokenType = AND
	case "class":
		tokenType = CLASS
	case "else":
		tokenType = ELSE
	case "false":
		tokenType = FALSE
	case "for":
		tokenType = FOR
	case "fun":
		tokenType = FUN
	case "if":
		tokenType = IF
	case "nil":
		tokenType = NIL
	case "or":
		tokenType = OR
	case "print":
		tokenType = PRINT
	case "return":
		tokenType = RETURN
	case "super":
		tokenType = SUPER
	case "this":
		tokenType = THIS
	case "true":
		tokenType = TRUE
	case "var":
		tokenType = VAR
	case "while":
		tokenType = WHILE
	default:
		tokenType = IDENTIFIER
	}

	s.addToken(tokenType, text)
}

func (s *Scanner) addToken(tokenType TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, text, literal, s.line))
}

func (s *Scanner) scanToken() ScannerError {
	c := s.advance()

	var err ScannerError
	switch c {
	case '(':
		s.addToken(LEFT_PAREN, nil)
	case ')':
		s.addToken(RIGHT_PAREN, nil)
	case '{':
		s.addToken(LEFT_BRACE, nil)
	case '}':
		s.addToken(RIGHT_BRACE, nil)
	case ',':
		s.addToken(COMMA, nil)
	case '.':
		s.addToken(DOT, nil)
	case '-':
		s.addToken(MINUS, nil)
	case '+':
		s.addToken(PLUS, nil)
	case ';':
		s.addToken(SEMICOLON, nil)
	case '*':
		s.addToken(STAR, nil)
	case '!':
		s.addToken(s.matchOrElse('=', BANG_EQUAL, BANG), nil)
	case '=':
		s.addToken(s.matchOrElse('=', EQUAL_EQUAL, EQUAL), nil)
	case '<':
		s.addToken(s.matchOrElse('=', LESS_EQUAL, LESS), nil)
	case '>':
		s.addToken(s.matchOrElse('=', GREATER_EQUAL, GREATER), nil)
	case '/':
		if s.match('/') {
			// A comment goes until the end of the Line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH, nil)
		}
	case '"':
		err = s.handleString()
	case '\r':
	case '\t':
	case ' ':
	case '\n':
		s.line += 1
	default:
		if s.isDigit(c) {
			err = s.handleNumber()
		} else if s.isAlpha(c) {
			s.handleIdentifier()
		} else {
			err = NewScannerError(s.current, "Unexpected character.")
		}
	}

	return err
}

func (s *Scanner) ScanTokens() ([]Token, ScannerError) {
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()

		if err != nil {
			return nil, err
		}
	}

	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))

	return s.tokens, nil
}
