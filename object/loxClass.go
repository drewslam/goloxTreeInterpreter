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

func (l *LoxClass) FindMethod(name string) (*LoxFunction, bool) {
	fmt.Printf("Looking for method: %s in class %s\n", name, l.Name)
	fmt.Printf("Available methods: %+v\n", l.Methods)

	if value, ok := l.Methods[name]; ok {
		fmt.Println("Method found:", name)
		return value, true
	}
	fmt.Println("Method not found", name)
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
	fmt.Printf("Instance created: %v\n", instance)

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
