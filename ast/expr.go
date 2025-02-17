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
	VisitCallExpr(expr *Call) interface{}
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

func (expr *Assign) Accept(visitor ExprVisitor) interface{} {
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

func (expr *Binary) Accept(visitor ExprVisitor) interface{} {
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

func (expr *Call) Accept(visitor ExprVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitCallExpr(expr)
}

// Grouping: Grouping Expression: "(expression)"
type Grouping struct {
	Expression Expr
}

func (expr *Grouping) Accept(visitor ExprVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitGroupingExpr(expr)
}

// Literal: Literal value: Number, String, true, false, nil
type Literal struct {
	Value interface{}
}

func (expr *Literal) Accept(visitor ExprVisitor) interface{} {
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

func (expr *Logical) Accept(visitor ExprVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitLogicalExpr(expr)
}

// Unary: Unary expression: "operator expression"
type Unary struct {
	Operator token.Token
	Right    Expr
}

func (expr *Unary) Accept(visitor ExprVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitUnaryExpr(expr)
}

// Variable expressions
type Variable struct {
	Name token.Token
}

func (expr *Variable) Accept(visitor ExprVisitor) interface{} {
	if visitor == nil {
		return nil
	}
	return visitor.VisitVariableExpr(expr)
}
