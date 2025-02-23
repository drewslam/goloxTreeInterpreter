package loxCallable

import (
	"time"

	"github.com/drewslam/goloxTreeInterpreter/environment"
)

func RegisterNatives(env *environment.Environment) {
	env.Define("clock", NewClockCallable())
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
