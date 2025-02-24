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
		return err
	}
	if err := l.run(string(bytes)); err != nil {
		return err
	}

	if errors.HadError {
		os.Exit(65)
	}

	return nil
}

func (l *Lox) runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		if line == "\n" {
			continue
		}
		if err := l.run(line); err != nil {
			fmt.Printf("Error executing line: %v\n", err)
		}
		errors.HadError = false
	}
}

func (l *Lox) run(source string) error {
	scanner := scanner.NewScanner(source)
	tokens := scanner.ScanTokens()

	parser := &parser.Parser{}
	parser.NewParser(tokens)
	statements := parser.Parse()

	// Stop if there is a syntax error
	if errors.HadError {
		os.Exit(65)
	}

	if errors.HadRuntimeError {
		os.Exit(70)
	}

	resolver := &resolver.Resolver{
		Interpreter:     l.interpreter,
		CurrentFunction: resolver.NOT_FUNCTION,
	}
	resolver.Resolve(statements)

	if errors.HadError {
		return nil
	}

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
			fmt.Printf("Error running file: %v\n", err)
			os.Exit(64)
		}
	default:
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	}
}
