package errors

import (
	"fmt"

	"github.com/drewslam/goloxTreeInterpreter/token"
)

// HadError tracks whether an error has occured
var HadError bool = false

// ParseError is a custom error type for parser errors
type ParseError struct {
	Message string
}

func (e *ParseError) Error() string {
	return e.Message
}

// NewParseError creates a new ParseError instance
func NewParseError(message string) *ParseError {
	return &ParseError{Message: message}
}

// ParseError function to report an error for a specific token and panic with a ParseError
func Error(token token.Token, message string) {
	if token.Type.String() == "EOF" {
		report(token.Line, " at end", message)
	} else {
		report(token.Line, " at '"+token.Lexeme+"'", message)
	}
	panic(NewParseError(message))
}

// Reports an error at a specific line
func ReportError(line int, message string) {
	report(line, "", message)
}

// Report formats and prints the error message
func report(line int, where string, message string) {
	fmt.Printf("[line %d] Error %v: %v\n", line, where, message)
	HadError = true
}
