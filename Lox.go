package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/drewslam/goloxTreeInterpreter/errors"
	"github.com/drewslam/goloxTreeInterpreter/interpreter"
	"github.com/drewslam/goloxTreeInterpreter/parser"
	"github.com/drewslam/goloxTreeInterpreter/resolver"
	"github.com/drewslam/goloxTreeInterpreter/scanner"
)

type Lox struct {
	interpreter *interpreter.Interpreter
}

func NewLox() *Lox {
	return &Lox{
		interpreter: interpreter.NewInterpreter(),
	}
}

func (l *Lox) runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Failed to read file: ", err)
	}
	return l.run(string(bytes))
}

func (l *Lox) runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Errorf("Error reading input: ", err)
			return
		}
		if line == "\n" {
			continue
		}
		if err := l.run(line); err != nil {
			fmt.Errorf("Error executing line: %v\n", err)
		}
	}
}

func (l *Lox) run(source string) error {
	scanner := scanner.NewScanner(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		return nil
	}

	parser := parser.NewParser(tokens)
	statements, err := parser.Parse()
	if err != nil {
		return nil
	}

	resolver := &resolver.Resolver{
		Interpreter:     l.interpreter,
		CurrentFunction: resolver.NOT_FUNCTION,
	}
	resolver.Resolve(statements)

	l.interpreter.Interpret(statements)
	return nil
}

func main() {
	lox := NewLox()

	switch len(os.Args) {
	case 1:
		lox.runPrompt()
	case 2:
		if err := lox.runFile(os.Args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)

			// Determine exit code based on error type
			var loxErr *errors.LoxError
			if errors.As(err, &loxErr) && loxErr.IsFatal {
				os.Exit(70) // Runtime error
			} else {
				os.Exit(65) // Syntax error
			}
		}
	default:
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	}
}
