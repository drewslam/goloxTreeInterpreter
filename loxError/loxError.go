package loxError

import (
	"fmt"

	"github.com/drewslam/goloxTreeInterpreter/token"
)

// LoxError represents an error in the interpreter
type LoxError struct {
	Line    int
	Where   string
	Message string
	IsFatal bool
}

// Error implements the error interface for RuntimeError.
func (e *LoxError) Error() string {
	return fmt.Sprintf("[line %d] Error %s: %s", e.Line, e.Where, e.Message)
}

// NewParseError creates a parse error (non-fatal)
func NewParseError(token token.Token, message string) *LoxError {
	where := " at end"
	if token.Type.String() != "EOF" {
		where = " at '" + token.Lexeme + "'"
	}
	return &LoxError{
		Line:    token.Line,
		Where:   where,
		Message: message,
		IsFatal: false,
	}
}

// NewRuntimeError creates a runtime error (fatal)
func NewRuntimeError(token token.Token, where string, message string) *LoxError {
	return &LoxError{
		Line:    token.Line,
		Where:   where,
		Message: message,
		IsFatal: true,
	}
}

// NewScanError creates a scan error (non-fatal)
func NewScanError(line int, message string) *LoxError {
	return &LoxError{
		Line:    line,
		Where:   "",
		Message: message,
		IsFatal: false,
	}
}

// Report error prints an error without panicking
func ReportError(err *LoxError) {
	fmt.Println(err.Error())
}

// ReportAndPanic reports an error and panics if it's fatal
func ReportAndPanic(err *LoxError) {
	fmt.Println(err.Error())
	if err.IsFatal {
		panic(err)
	}
}
