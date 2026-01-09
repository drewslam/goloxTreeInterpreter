package object

import (
	"github.com/drewslam/goloxTreeInterpreter/loxCallable"
)

// type Interpreter any

type LoxClass struct {
	Name       string
	Superclass *LoxClass
	Methods    map[string]*LoxFunction
}

func (l *LoxClass) FindMethod(name string) (*LoxFunction, bool) {
	if value, ok := l.Methods[name]; ok {
		return value, true
	}

	if l.Superclass != nil {
		return l.Superclass.FindMethod(name)
	}

	return nil, false
}

func (l *LoxClass) String() string {
	return l.Name
}

func (l *LoxClass) Call(interpreter loxCallable.Interpreter, arguments []any) any {
	instance := &LoxInstance{
		Klass:  l,
		Fields: make(map[string]any),
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
