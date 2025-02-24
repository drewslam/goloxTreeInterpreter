package object

import (
	"fmt"

	"github.com/drewslam/goloxTreeInterpreter/errors"
	"github.com/drewslam/goloxTreeInterpreter/token"
)

type LoxInstance struct {
	Klass  *LoxClass
	Fields map[string]interface{}
}

func (l *LoxInstance) String() string {
	return l.Klass.Name + " instance"
}

func (l *LoxInstance) Get(name token.Token) interface{} {
	/*if value, exists := l.Fields[name.Lexeme]; exists {
		return value
	}*/

	method, exists := l.Klass.FindMethod(name.Lexeme)
	if exists {
		fmt.Printf("Binding method: %v\n", method)
		return method.Bind(l)
	}

	return errors.NewRuntimeError(name, "Undefined property '"+name.Lexeme+"'.")
}

func (l *LoxInstance) Set(name token.Token, value interface{}) {
	l.Fields[name.Lexeme] = value
}
