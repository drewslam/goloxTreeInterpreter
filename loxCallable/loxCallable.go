package loxCallable

import (
	"time"

	"github.com/drewslam/goloxTreeInterpreter/ast"
	"github.com/drewslam/goloxTreeInterpreter/environment"
)

type Interpreter interface {
	ExecuteBlock(statements []ast.Stmt, environment *environment.Environment)
	GetGlobals() *environment.Environment
}

type LoxCallable interface {
	Arity() int
	Call(interpreter Interpreter, arguments []interface{}) interface{}
	String() string
}

func NewClockCallable() LoxCallable {
	return &Clock{}
}

type Clock struct{}

func (c *Clock) Arity() int {
	return 0
}

func (c *Clock) Call(interpreter Interpreter, arguments []interface{}) interface{} {
	return float64(time.Now().UnixNano()) / 1e9
}

func (c *Clock) String() string {
	return "<native fn>"
}

var _ LoxCallable = (*Clock)(nil)
