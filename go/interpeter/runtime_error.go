package interpeter

import "fmt"

type RuntimeError interface {
	error
	Error() string
	Line() int
	Message() string
}

type BaseRuntimeError struct {
	RuntimeError

	line    int
	message string
}

func NewRuntimeError(line int, message string) RuntimeError {
	return &BaseRuntimeError{
		line: line, message: message,
	}
}

func (e *BaseRuntimeError) Error() string {
	return fmt.Sprintf("[line %d] RuntimeError: %v", e.Line(), e.Message())
}

func (e *BaseRuntimeError) Line() int {
	return e.line
}

func (e *BaseRuntimeError) Message() string {
	return e.message
}
