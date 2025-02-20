package object

import (
	"github.com/drewslam/goloxTreeInterpreter/loxCallable"
)

type Interpreter interface{}

type LoxClass struct {
	loxCallable loxCallable.LoxCallable
	Name        string
}

func (l *LoxClass) ToString() string {
	return l.Name
}

func (l *LoxClass) Call(interpreter Interpreter, arguments []interface{}) interface{} {
	instance := &LoxInstance{
		Klass: *l,
	}
	return instance
}

func (l *LoxClass) Arity() int {
	return 0
}
