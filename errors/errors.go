package errors

import "fmt"

var HadError bool = false

func ReportError(line int, message string) {
	Report(line, "", message)
}

func Report(line int, where string, message string) {
	fmt.Printf("[line %d] Error %v: %v\n", line, where, message)
	HadError = true
}
