package object

import (
	"github.com/drewslam/goloxTreeInterpreter/errors"
	"github.com/drewslam/goloxTreeInterpreter/token"
)

type LoxInstance struct {
	Klass  LoxClass
	Fields map[string]interface{}
}

func (l *LoxInstance) ToString() string {
	return l.Klass.Name + " instance"
}

func (l *LoxInstance) Get(name token.Token) interface{} {
	if value, exists := l.Fields[name.Lexeme]; exists {
		return value
	}
	return errors.NewRuntimeError(name, "Undefined property '"+name.Lexeme+"'.")
}

func (l *LoxInstance) Set(name token.Token, value interface{}) {
	l.Fields[name.Lexeme] = value
}
