package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/drewslam/goloxTreeInterpreter/errors"
	"github.com/drewslam/goloxTreeInterpreter/parser"
	"github.com/drewslam/goloxTreeInterpreter/scanner"
)

func runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	run(string(bytes))
	if errors.HadError {
		os.Exit(65)
	}
	return nil
}

func runPrompt() {
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
		run(line)
		errors.HadError = false
	}
}

func run(source string) {
	sc := scanner.NewScanner(source)
	tokens := sc.ScanTokens()
	pa := parser.Parser{}
	pa.NewParser(tokens)
	expr := pa.Parse()

	// Stop if there is a syntax error
	if errors.HadError {
		return
	}

	// For now, we'll print the parsed expression
	fmt.Println(expr)
}

func main() {
	argc := len(os.Args)

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
}
