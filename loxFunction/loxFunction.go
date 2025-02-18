package loxFunction

import (
	"fmt"

	"github.com/drewslam/goloxTreeInterpreter/ast"
	"github.com/drewslam/goloxTreeInterpreter/environment"
	"github.com/drewslam/goloxTreeInterpreter/loxCallable"
	"github.com/drewslam/goloxTreeInterpreter/returnValue"
)

/*type Interpreter interface {
	ExecuteBlock(statements []ast.Stmt, environment *environment.Environment)
	GetGlobals() *environment.Environment
}*/

/*type LoxFunction interface {
	Arity() int
	Call(interpreter Interpreter, arguments []interface{}) interface{}
	String() string
}*/

type loxFunction struct {
	declaration *ast.Function
}

func NewLoxFunction(declaration *ast.Function) *loxFunction {
	return &loxFunction{
		declaration: declaration,
	}
}

func (l *loxFunction) String() string {
	return fmt.Sprintf("<fn %v>", l.declaration.Name.Lexeme)
}

func (l *loxFunction) Arity() int {
	return len(l.declaration.Params)
}

func (l *loxFunction) Call(interpreter loxCallable.Interpreter, arguments []interface{}) (result interface{}) {
	environment := environment.NewEnvironment(interpreter.GetGlobals())

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
