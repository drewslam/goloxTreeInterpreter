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
	//	VisitLogicalExpr(expr *Logical) interface{}
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
		panic("Visitor is nil in Assign.Accept")
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
		panic("Visitor is nil in Binary.Accept")
	}
	return visitor.VisitBinaryExpr(b)
}

// Grouping: Grouping Expression: "(expression)"
type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(visitor ExprVisitor) interface{} {
	if visitor == nil {
		panic("Visitor is nil in Grouping.Accept")
	}
	return visitor.VisitGroupingExpr(g)
}

// Literal: Literal value: Number, String, true, false, nil
type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(visitor ExprVisitor) interface{} {
	if visitor == nil {
		panic("Visitor is nil in Literal.Accept")
	}
	return visitor.VisitLiteralExpr(l)
}

// Unary: Unary expression: "operator expression"
type Unary struct {
	Operator token.Token
	Right    Expr
}

func (u *Unary) Accept(visitor ExprVisitor) interface{} {
	if visitor == nil {
		panic("Visitor is nil in Unary.Accept")
	}
	return visitor.VisitUnaryExpr(u)
}

// Variable expressions
type Variable struct {
	Name token.Token
}

func (v *Variable) Accept(visitor ExprVisitor) interface{} {
	if visitor == nil {
		panic("Visitor is nil in Variable.Accept")
	}
	return visitor.VisitVariableExpr(v)
}
