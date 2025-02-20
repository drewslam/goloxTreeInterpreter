package object

import (
	"fmt"

	"github.com/drewslam/goloxTreeInterpreter/loxCallable"
)

// type Interpreter interface{}

type LoxClass struct {
	Name    string
	Methods map[string]*LoxFunction
}

func (l *LoxClass) FindMethod(name string) *LoxFunction {
	if value, ok := l.Methods[name]; ok {
		return value
	}
	return nil
}

func (l *LoxClass) String() string {
	return l.Name
}

func (l *LoxClass) Call(interpreter loxCallable.Interpreter, arguments []interface{}) interface{} {
	instance := &LoxInstance{
		Klass:  l,
		Fields: make(map[string]interface{}),
	}
	fmt.Printf("Instance created: %v\n", instance)

	for name, method := range l.Methods {
		boundMethod := method.Bind(instance)
		instance.Fields[name] = boundMethod
		fmt.Printf("Method: %s\n", instance.Fields[name])
	}
	return instance
}

func (l *LoxClass) Arity() int {
	return 0
}

var _ loxCallable.LoxCallable = (*LoxClass)(nil)
