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

func (e *Environment) Get(name token.Token) (interface{}, *errors.LoxError) {
	if value, exists := e.Values[name.Lexeme]; exists {
		return value, nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	errMsg := "Undefined variable: '" + name.Lexeme + "'"

	return nil, errors.NewRuntimeError(name, fmt.Sprintf("%d", name.Line), errMsg)
}

func (e *Environment) Assign(name token.Token, value interface{}) *errors.LoxError {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.Assign(name, value)
	}

	errMsg := "Undefined variable: '" + name.Lexeme + "'"

	return errors.NewRuntimeError(name, fmt.Sprintf("%d", name.Line), errMsg)
}

func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Environment) Ancestor(distance int) (*Environment, error) {
	environment := e
	for i := 0; i < distance; i++ {
		if environment.Enclosing == nil {
			return nil, fmt.Errorf("Environment chain not deep enough")
		}
		environment = environment.Enclosing
	}
	return environment, nil
}

func (e *Environment) GetAt(distance int, name string) (interface{}, error) {
	if res, err := e.Ancestor(distance).Values[name]; err != nil {
		return res, nil
	} else {
		return nil, err
	}
}

func (e *Environment) AssignAt(distance int, name token.Token, value interface{}) {
	e.Ancestor(distance).Values[name.Lexeme] = value
}
