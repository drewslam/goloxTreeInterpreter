package object

import (
	"github.com/drewslam/goloxTreeInterpreter/loxCallable"
)

// type Interpreter interface{}

type LoxClass struct {
	Name    string
	Methods map[string]*LoxFunction
}

func (l *LoxClass) FindMethod(name string) (*LoxFunction, bool) {

	if value, ok := l.Methods[name]; ok {
		return value, true
	}

	return nil, false
}

func (l *LoxClass) String() string {
	return l.Name
}

func (l *LoxClass) Call(interpreter loxCallable.Interpreter, arguments []interface{}) interface{} {
	instance := &LoxInstance{
		Klass:  l,
		Fields: make(map[string]interface{}),
	}

	initializer, exists := l.FindMethod("init")
	if exists {
		initializer.Bind(instance).Call(interpreter, arguments)
	}
	return instance
}

func (l *LoxClass) Arity() int {
	initializer, exists := l.FindMethod("init")
	if !exists {
		return 0
	}
	return initializer.Arity()
}

var _ loxCallable.LoxCallable = (*LoxClass)(nil)
