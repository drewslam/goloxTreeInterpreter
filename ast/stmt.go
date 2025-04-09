package ast

import (
	"github.com/drewslam/goloxTreeInterpreter/token"
)

type Stmt interface {
	Accept(visitor StmtVisitor) interface{}
}

type StmtVisitor interface {
	VisitBlockStmt(stmt *Block) interface{}
	VisitClassStmt(stmt *Class) interface{}
	VisitExpressionStmt(stmt *Expression) interface{}
	VisitFunctionStmt(stmt *Function) interface{}
	VisitIfStmt(stmt *If) interface{}
	VisitPrintStmt(stmt *Print) interface{}
	VisitReturnStmt(stmt *Return) interface{}
	VisitVarStmt(stmt *Var) interface{}
	VisitWhileStmt(stmt *While) interface{}
}

// Block type
type Block struct {
	Statements []Stmt
}

func (stmt *Block) Accept(visitor StmtVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitBlockStmt(stmt)
}

func NewBlockStmt(statements []Stmt) *Block {
	return &Block{Statements: statements}
}

// Class type
type Class struct {
	Name       token.Token
	Superclass *Variable
	Methods    []*Function
}

func (stmt *Class) Accept(visitor StmtVisitor) interface{} {
	if visitor == nil {
		return nil
	}

	return visitor.VisitClassStmt(stmt)
}

func NewClassStmt(name token.Token, superclass *Variable, methods []*Function) *Class {
	return &Class{
		Name:       name,
		Superclass: superclass,
		Methods:    methods,
	}
}

// Expression type
type Expression struct {
	Expr Expr
}

func (stmt *Expression) Accept(visitor StmtVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitExpressionStmt(stmt)
}

func NewExpressionStmt(expr Expr) *Expression {
	return &Expression{Expr: expr}
}

// Function type
type Function struct {
	Name   token.Token
	Params []token.Token
	Body   []Stmt
}

func (stmt *Function) Accept(visitor StmtVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitFunctionStmt(stmt)
}

func NewFunctionStmt(name token.Token, params []token.Token, body []Stmt) *Function {
	return &Function{
		Name:   name,
		Params: params,
		Body:   body,
	}
}

// If type
type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (stmt *If) Accept(visitor StmtVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitIfStmt(stmt)
}

func NewIfStmt(condition Expr, thenBranch Stmt, elseBranch Stmt) *If {
	return &If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

// Print type
type Print struct {
	Expr Expr
}

func (stmt *Print) Accept(visitor StmtVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitPrintStmt(stmt)
}

func NewPrintStmt(expr Expr) *Print {
	return &Print{Expr: expr}
}

// Return type
type Return struct {
	Keyword token.Token
	Value   Expr
}

func (stmt *Return) Accept(visitor StmtVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitReturnStmt(stmt)
}

// Variable type
type Var struct {
	Name        token.Token
	Initializer Expr
}

func (stmt *Var) Accept(visitor StmtVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitVarStmt(stmt)
}

func NewVarStmt(name token.Token, initializer Expr) *Var {
	return &Var{
		Name:        name,
		Initializer: initializer,
	}
}

// While type
type While struct {
	Condition Expr
	Body      Stmt
}

func (stmt *While) Accept(visitor StmtVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitWhileStmt(stmt)
}
