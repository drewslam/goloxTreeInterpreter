package loxCallable

import (
	"github.com/drewslam/goloxTreeInterpreter/ast"
	"github.com/drewslam/goloxTreeInterpreter/environment"
)

type Interpreter interface {
	ExecuteBlock(statements []ast.Stmt, environment *environment.Environment) any
	GetGlobals() *environment.Environment
}

type LoxCallable interface {
	Arity() int
	Call(interpreter Interpreter, arguments []any) any
	String() string
}

func NewClockCallable() LoxCallable {
	return &Clock{}
}

func RegisterNatives(env *environment.Environment) {
	env.Define("clock", NewClockCallable())
}
