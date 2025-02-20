package object

import (
	"fmt"

	"github.com/drewslam/goloxTreeInterpreter/ast"
	"github.com/drewslam/goloxTreeInterpreter/environment"
	"github.com/drewslam/goloxTreeInterpreter/loxCallable"
	"github.com/drewslam/goloxTreeInterpreter/returnValue"
)

type LoxFunction struct {
	Declaration *ast.Function
	Closure     *environment.Environment
}

func NewLoxFunction(declaration *ast.Function, closure *environment.Environment) *LoxFunction {
	return &LoxFunction{
		Declaration: declaration,
		Closure:     closure,
	}
}

func (l *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	environment := &environment.Environment{
		Enclosing: l.Closure,
	}
	environment.Define("this", instance)
	return &LoxFunction{
		Declaration: l.Declaration,
		Closure:     environment,
	}
}

func (l *LoxFunction) String() string {
	return fmt.Sprintf("<fn %v>", l.Declaration.Name.Lexeme)
}

func (l *LoxFunction) Arity() int {
	return len(l.Declaration.Params)
}

func (l *LoxFunction) Call(interpreter loxCallable.Interpreter, arguments []interface{}) (result interface{}) {
	environment := environment.NewEnvironment(l.Closure)

	for i, param := range l.Declaration.Params {
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

	interpreter.ExecuteBlock(l.Declaration.Body, environment)

	return
}

var _ loxCallable.LoxCallable = (*LoxFunction)(nil)
