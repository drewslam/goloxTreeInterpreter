package ast

import (
	"github.com/drewslam/goloxTreeInterpreter/token"
)

type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}

type ExprVisitor interface {
	VisitAssignExpr(expr *Assign) interface{}
	VisitBinaryExpr(expr *Binary) interface{}
	//	VisitCallExpr(expr *Call) interface{}
	//	VisitGetExpr(expr *Get) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitLogicalExpr(expr *Logical) interface{}
	//	VisitSetExpr(expr *Set) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
	VisitVariableExpr(expr *Variable) interface{}
}

// Assignment
type Assign struct {
	Name  token.Token
	Value Expr
}

func (a *Assign) Accept(visitor ExprVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitAssignExpr(a)
}

// Binary: Binary Expression: "left operator right"
type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b *Binary) Accept(visitor ExprVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitBinaryExpr(b)
}

// Grouping: Grouping Expression: "(expression)"
type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(visitor ExprVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitGroupingExpr(g)
}

// Literal: Literal value: Number, String, true, false, nil
type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(visitor ExprVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitLiteralExpr(l)
}

// Logical expressions:
type Logical struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (l *Logical) Accept(visitor ExprVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitLogicalExpr(l)
}

// Unary: Unary expression: "operator expression"
type Unary struct {
	Operator token.Token
	Right    Expr
}

func (u *Unary) Accept(visitor ExprVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitUnaryExpr(u)
}

// Variable expressions
type Variable struct {
	Name token.Token
}

func (v *Variable) Accept(visitor ExprVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitVariableExpr(v)
}
