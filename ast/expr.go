package ast

import "github.com/drewslam/goloxTreeInterpreter/token"

type Expr interface {
	Accept(visitor ExprVisitor) any
}

type ExprVisitor interface {
	VisitAssignExpr(expr *Assign) any
	VisitBinaryExpr(expr *Binary) any
	VisitCallExpr(expr *Call) any
	VisitGetExpr(expr *Get) any
	VisitGroupingExpr(expr *Grouping) any
	VisitLiteralExpr(expr *Literal) any
	VisitLogicalExpr(expr *Logical) any
	VisitSetExpr(expr *Set) any
	VisitSuperExpr(expr *Super) any
	VisitThisExpr(expr *This) any
	VisitUnaryExpr(expr *Unary) any
	VisitVariableExpr(expr *Variable) any
}

// Assignment
type Assign struct {
	Name  token.Token
	Value Expr
}

func (expr *Assign) Accept(visitor ExprVisitor) any {
	if visitor == nil {
		return nil
	}
	return visitor.VisitAssignExpr(expr)
}

// Binary: Binary Expression: "left operator right"
type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (expr *Binary) Accept(visitor ExprVisitor) any {
	if visitor == nil {
		return nil
	}
	return visitor.VisitBinaryExpr(expr)
}

// Call
type Call struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
}

func (expr *Call) Accept(visitor ExprVisitor) any {
	if visitor == nil {
		return nil
	}
	return visitor.VisitCallExpr(expr)
}

// Get
type Get struct {
	Object Expr
	Name   token.Token
}

func (expr *Get) Accept(visitor ExprVisitor) any {
	if visitor == nil {
		return nil
	}
	return visitor.VisitGetExpr(expr)
}

// Grouping: Grouping Expression: "(expression)"
type Grouping struct {
	Expression Expr
}

func (expr *Grouping) Accept(visitor ExprVisitor) any {
	if visitor == nil {
		return nil
	}
	return visitor.VisitGroupingExpr(expr)
}

// Literal: Literal value: Number, String, true, false, nil
type Literal struct {
	Value any
}

func (expr *Literal) Accept(visitor ExprVisitor) any {
	if visitor == nil {
		return nil
	}
	return visitor.VisitLiteralExpr(expr)
}

// Logical expressions:
type Logical struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (expr *Logical) Accept(visitor ExprVisitor) any {
	if visitor == nil {
		return nil
	}
	return visitor.VisitLogicalExpr(expr)
}

// Set
type Set struct {
	Object Expr
	Name   token.Token
	Value  Expr
}

func (expr *Set) Accept(visitor ExprVisitor) any {
	if visitor == nil {
		return nil
	}
	return visitor.VisitSetExpr(expr)
}

// Super
type Super struct {
	Keyword token.Token
	Method  token.Token
}

func (expr *Super) Accept(visitor ExprVisitor) any {
	if visitor == nil {
		return nil
	}
	return visitor.VisitSuperExpr(expr)
}

// This
type This struct {
	Keyword token.Token
}

func (expr *This) Accept(visitor ExprVisitor) any {
	if visitor == nil {
		return nil
	}
	return visitor.VisitThisExpr(expr)
}

// Unary: Unary expression: "operator expression"
type Unary struct {
	Operator token.Token
	Right    Expr
}

func (expr *Unary) Accept(visitor ExprVisitor) any {
	if visitor == nil {
		return nil
	}
	return visitor.VisitUnaryExpr(expr)
}

// Variable expressions
type Variable struct {
	Name token.Token
}

func (expr *Variable) Accept(visitor ExprVisitor) any {
	if visitor == nil {
		return nil
	}
	return visitor.VisitVariableExpr(expr)
}
