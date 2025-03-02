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

func (l *LoxFunction) Call(interpreter loxCallable.Interpreter, arguments []interface{}) interface{} {
	environment := environment.NewEnvironment(l.Closure)

	// Ensure argument count matches parameter count
	if len(arguments) != len(l.Declaration.Params) {
		panic(fmt.Errorf("Expected %d arguments but got %d.", len(l.Declaration.Params), len(arguments)))
	}

	// Define function parameters in environment
	for i, param := range l.Declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	// Try executing the function
	defer func() {
		if r := recover(); r != nil {
			if rv, ok := r.(*returnValue.ReturnValue); ok {
				// Set the result via the deferred function
				if l.IsInitializer {
					thisVal, _ := l.Closure.GetAt(0, "this")
					panic(&returnValue.ReturnValue{Value: thisVal})
				}
				panic(rv) // Rethrow the return value to the caller
			}
			panic(r) // Re-panic other errors
		}
	}()

	interpreter.ExecuteBlock(l.Declaration.Body, environment)

	// Ensure constructors return 'this'
	if l.IsInitializer {
		val, _ := l.Closure.GetAt(0, "this")
		return val
	}

	return nil
}

var _ loxCallable.LoxCallable = (*LoxFunction)(nil)
