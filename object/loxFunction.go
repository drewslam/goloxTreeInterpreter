package object

import (
	"fmt"

	"github.com/drewslam/goloxTreeInterpreter/ast"
	"github.com/drewslam/goloxTreeInterpreter/environment"
	"github.com/drewslam/goloxTreeInterpreter/loxCallable"
	"github.com/drewslam/goloxTreeInterpreter/returnValue"
)

type LoxFunction struct {
	Declaration   *ast.Function
	Closure       *environment.Environment
	IsInitializer bool
}

func NewLoxFunction(declaration *ast.Function, closure *environment.Environment, isInitializer bool) *LoxFunction {
	return &LoxFunction{
		IsInitializer: isInitializer,
		Declaration:   declaration,
		Closure:       closure,
	}
}

func (l *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	environment := &environment.Environment{
		Enclosing: l.Closure,
		Values:    make(map[string]interface{}),
	}
	environment.Define("this", instance)
	return &LoxFunction{
		Declaration:   l.Declaration,
		Closure:       environment,
		IsInitializer: l.IsInitializer,
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

	// var result interface{} = nil

	defer func() {
		if r := recover(); r != nil {
			if returnValue, ok := r.(*returnValue.ReturnValue); ok {
				if l.IsInitializer {
					result = l.Closure.GetAt(0, "this")
				} else {
					result = returnValue.Value
				}
				return
			} else {
				panic(r)
			}
		}
	}()

	interpreter.ExecuteBlock(l.Declaration.Body, environment)

	if result == nil && l.IsInitializer {
		result = l.Closure.GetAt(0, "this")
	}

	return //result
}

var _ loxCallable.LoxCallable = (*LoxFunction)(nil)
