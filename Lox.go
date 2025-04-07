package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/drewslam/goloxTreeInterpreter/interpreter"
	"github.com/drewslam/goloxTreeInterpreter/loxDebug"
	"github.com/drewslam/goloxTreeInterpreter/loxError"
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
		return fmt.Errorf("Failed to read file: %v", err)
	}
	result := l.run(string(bytes))
	if result == nil {
		return nil
	}
	return result
}

func (l *Lox) runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			return
		}
		if line == "\n" {
			continue
		}
		if err := l.run(line); err != nil {
			fmt.Fprintf(os.Stderr, "Error executing line: %v\n", err)
		}
	}
}

func (l *Lox) run(source string) *loxError.LoxError {
	scanner := scanner.NewScanner(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		return err
	}

	parser := parser.NewParser(tokens)
	statements, err := parser.Parse()
	if err != nil {
		return err
	}

	resolver := resolver.NewResolver(l.interpreter)
	err = resolver.Resolve(statements)
	if err != nil {
		return err
	}

	l.interpreter.Interpret(statements)
	return nil
}

func main() {
	loxDebug.InitializeLogger()
	defer loxDebug.CloseLogger()

	lox := NewLox()

	switch len(os.Args) {
	case 1:
		lox.runPrompt()
	case 2:
		err := lox.runFile(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)

			// Determine exit code based on error type
			var loxErr *loxError.LoxError
			if errors.As(err, &loxErr) {
				if loxErr.IsFatal {
					os.Exit(70) // Runtime error
				} else {
					os.Exit(65) // Syntax error
				}
			} else {
				os.Exit(65)
			}
		}
	default:
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	}
}
