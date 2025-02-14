package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/drewslam/goloxTreeInterpreter/errors"
	"github.com/drewslam/goloxTreeInterpreter/interpreter"
	"github.com/drewslam/goloxTreeInterpreter/parser"
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
	expr := parser.Parse()

	// Stop if there is a syntax error
	if errors.HadError {
		os.Exit(65)
	}

	if errors.HadRuntimeError {
		os.Exit(70)
	}

	l.interpreter.Interpret(expr)
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
	/*
		if argc > 2 {
			fmt.Println("Usage: golox [script]")
			os.Exit(64)
		} else if argc == 2 {

			err := runFile(os.Args[1])
			if err != nil {
				fmt.Println("Error running file: ", err)
				os.Exit(64)
			}
		} else {

			runPrompt()
		}
	*/
}
