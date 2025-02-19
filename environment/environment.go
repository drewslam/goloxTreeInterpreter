package environment

import (
	"fmt"

	"github.com/drewslam/goloxTreeInterpreter/errors"
	"github.com/drewslam/goloxTreeInterpreter/token"
)

type Environment struct {
	Enclosing *Environment
	Values    map[string]interface{}
}

/*func NewEnvrionment() *Environment {
	return &Environment{
		Enclosing: nil,
		Values:    make(map[string]interface{}),
	}
}*/

func NewEnvironment(enclosing ...*Environment) *Environment {
	var parent *Environment
	if len(enclosing) > 0 {
		parent = enclosing[0]
	}

	return &Environment{
		Enclosing: parent,
		Values:    make(map[string]interface{}),
	}
}

func (e *Environment) Get(name token.Token) interface{} {
	if value, exists := e.Values[name.Lexeme]; exists {
		return value
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	panic(errors.NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'"))
}

func (e *Environment) Assign(name token.Token, value interface{}) error {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.Assign(name, value)
	}

	return fmt.Errorf("Undefined variable '%s'", name.Lexeme)
}

func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Environment) Ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.Enclosing
	}
	return environment
}

func (e *Environment) GetAt(distance int, name string) interface{} {
	return e.Ancestor(distance).Values[name]
}

func (e *Environment) AssignAt(distance int, name token.Token, value interface{}) {
	e.Ancestor(distance).Values[name.Lexeme] = value
}
