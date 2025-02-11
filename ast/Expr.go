package ast

import (
	"github.com/drewslam/goloxTreeInterpreter/token"
)

type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}

type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
}

// Binary: Binary Expression: "left operator right"
type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b *Binary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(b)
}

// Grouping: Grouping Expression: "(expression)"
type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGroupingExpr(g)
}

// Literal: Literal value: Number, String, true, false, nil
type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}

// Unary: Unary expression: "operator expression"
type Unary struct {
	Operator token.Token
	Right    Expr
}

func (u *Unary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}
