package ast

import "github.com/drewslam/goloxTreeInterpreter/token"

type Stmt interface {
	Accept(visitor StmtVisitor) interface{}
}

type StmtVisitor interface {
	VisitBlockStmt(stmt *Block) interface{}
	VisitExpressionStmt(stmt *Expression) interface{}
	VisitIfStmt(stmt *If) interface{}
	VisitPrintStmt(stmt *Print) interface{}
	VisitVarStmt(stmt *Var) interface{}
}

// Block type
type Block struct {
	Statements []Stmt
}

func (b *Block) Accept(visitor StmtVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitBlockStmt(b)
}

func NewBlockStmt(statements []Stmt) *Block {
	return &Block{Statements: statements}
}

// Expression type
type Expression struct {
	Expr Expr
}

func (e *Expression) Accept(visitor StmtVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitExpressionStmt(e)
}

func NewExpressionStmt(expr Expr) *Expression {
	return &Expression{Expr: expr}
}

// If type
type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i *If) Accept(visitor StmtVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitIfStmt(i)
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

func (p *Print) Accept(visitor StmtVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitPrintStmt(p)
}

func NewPrintStmt(expr Expr) *Print {
	return &Print{Expr: expr}
}

// Variable type
type Var struct {
	Name        token.Token
	Initializer Expr
}

func (v *Var) Accept(visitor StmtVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitVarStmt(v)
}

func NewVarStmt(name token.Token, initializer Expr) *Var {
	return &Var{
		Name:        name,
		Initializer: initializer,
	}
}
