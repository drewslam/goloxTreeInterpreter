package object

import (
	"fmt"

	"github.com/drewslam/goloxTreeInterpreter/ast"
	"github.com/drewslam/goloxTreeInterpreter/environment"
	"github.com/drewslam/goloxTreeInterpreter/loxCallable"
	"github.com/drewslam/goloxTreeInterpreter/loxError"
	"github.com/drewslam/goloxTreeInterpreter/returnValue"
	"github.com/drewslam/goloxTreeInterpreter/token"
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
	env := environment.NewEnvironment(l.Closure)
	/*environment := &environment.Environment{
		Enclosing: l.Closure,
		Values:    make(map[string]interface{}),
	}*/
	env.Define("this", instance)
	return &LoxFunction{
		Declaration:   l.Declaration,
		Closure:       env,
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
	// Ensure argument count matches parameter count
	if len(arguments) != len(l.Declaration.Params) {
		message := (fmt.Sprintf("Expected %d arguments but got %d.", len(l.Declaration.Params), len(arguments)))
		err := loxError.NewRuntimeError(l.Declaration.Name, "", message)
		loxError.ReportAndPanic(err)
	}

	env := environment.NewEnvironment(l.Closure)

	// Define function parameters in environment
	for i, param := range l.Declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	for _, param := range l.Declaration.Params {
		val, err := env.Get(token.Token{Lexeme: param.Lexeme})
		fmt.Printf("Parameter %s = %v (error %v)\n", param.Lexeme, val, err)
	}

	var result interface{} = nil

	// Try executing the function
	func() {
		defer func() {
			if r := recover(); r != nil {
				if rv, ok := r.(*returnValue.ReturnValue); ok {
					fmt.Println("Recovered return value:", rv.Value)
					result = rv.Value
					fmt.Println("Stored returnVal in Call():", result)
					// return
				} else {
					fmt.Println("Recovered unknown panic:", r)
					panic(r) // Re-panic other errors
				}
			}
		}()

		fmt.Println("Executing function body for:", l.String())
		interpreter.ExecuteBlock(l.Declaration.Body, env)
	}()

	// Ensure constructors return 'this'
	if l.IsInitializer {
		val, _ := l.Closure.GetAt(0, "this")
		fmt.Println("Returning 'this' from Call()")
		return val
	}

	fmt.Println("Returning from Call():", result)
	return result
}

var _ loxCallable.LoxCallable = (*LoxFunction)(nil)
