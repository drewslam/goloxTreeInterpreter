package errors

import (
	"fmt"

	"github.com/drewslam/goloxTreeInterpreter/token"
)

// HadError tracks whether a syntax error has occured
var HadError bool = false
var HadRuntimeError bool = false

// ParseError is a custom error type for parser errors.
type ParseError struct {
	Message string
}

// RuntimeError is a custom error type for runtime errors.
type RuntimeError struct {
	Token   token.Token
	Message string
}

// Error implements the error interface for ParseError.
func (e *ParseError) Error() string {
	return e.Message
}

// Error implements the error interface for RuntimeError.
func (r *RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] %s", r.Token.Line, r.Message)
}

// NewParseError creates a new ParseError instance and marks HadError as true.
func NewParseError(message string) *ParseError {
	HadError = true
	return &ParseError{Message: message}
}

// NewRuntimeError creates a new RuntimeError instance and marks HadError as true.
func NewRuntimeError(token token.Token, message string) *RuntimeError {
	HadRuntimeError = true
	return &RuntimeError{
		Token:   token,
		Message: message,
	}
}

// ParseError function to report an error for a specific token and panic with a ParseError
func ReportParseError(token token.Token, message string) {
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
