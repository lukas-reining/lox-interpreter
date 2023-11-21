package scanner

import "fmt"

type ScannerError interface {
	error
	Error() string
	Line() int
	Message() string
}

type BaseScannerError struct {
	ScannerError

	line    int
	message string
}

func NewScannerError(line int, message string) ScannerError {
	return &BaseScannerError{
		line: line, message: message,
	}
}

func (e *BaseScannerError) Error() string {
	return fmt.Sprintf("[line %d] RuntimeError: %v", e.Line(), e.Message())
}

func (e *BaseScannerError) Line() int {
	return e.line
}

func (e *BaseScannerError) Message() string {
	return e.message
}
