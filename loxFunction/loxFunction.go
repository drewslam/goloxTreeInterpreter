package loxFunction

import (
	"fmt"

	"github.com/drewslam/goloxTreeInterpreter/ast"
	"github.com/drewslam/goloxTreeInterpreter/environment"
	"github.com/drewslam/goloxTreeInterpreter/loxCallable"
	"github.com/drewslam/goloxTreeInterpreter/returnValue"
)

type loxFunction struct {
	declaration *ast.Function
	closure     *environment.Environment
}

func NewLoxFunction(declaration *ast.Function, closure *environment.Environment) *loxFunction {
	return &loxFunction{
		declaration: declaration,
		closure:     closure,
	}
}

func (l *loxFunction) String() string {
	return fmt.Sprintf("<fn %v>", l.declaration.Name.Lexeme)
}

func (l *loxFunction) Arity() int {
	return len(l.declaration.Params)
}

func (l *loxFunction) Call(interpreter loxCallable.Interpreter, arguments []interface{}) (result interface{}) {
	environment := environment.NewEnvironment(l.closure)

	for i, param := range l.declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	defer func() {
		if r := recover(); r != nil {
			if returnValue, ok := r.(*returnValue.ReturnValue); ok {
				result = returnValue.Value
			} else {
				panic(r)
			}
		}
	}()

	interpreter.ExecuteBlock(l.declaration.Body, environment)

	return
}

var _ loxCallable.LoxCallable = (*loxFunction)(nil)
