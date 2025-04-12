package loxError

import (
	"fmt"
	"os"

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
	return fmt.Sprintf("[line %d] Error at %s: %s", e.Line, e.Where, e.Message)
}

// NewParseError creates a parse error (non-fatal)
func NewParseError(token token.Token, message string) *LoxError {
	where := "end"
	if token.Type.String() != "EOF" {
		where = token.Lexeme
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
	panic(err)
}

func HandleRecoveredError(r any) {
	if err, ok := r.(*LoxError); ok {
		fmt.Println(err.Error())
		if err.IsFatal {
			os.Exit(1)
		}
		// Error is non-fatal
		return
	}

	// Not a LoxError -- Repanic
	panic(r)
}
