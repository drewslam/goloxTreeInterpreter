package ast

type Stmt interface {
	Accept(visitor StmtVisitor) interface{}
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt *Expression) interface{}
	VisitPrintStmt(stmt *Print) interface{}
}

// Expression type
type Expression struct {
	Expr Expr
}

func (e *Expression) Accept(visitor StmtVisitor) interface{} {
	if visitor == nil {
		panic("Visitor is nil in Expression.Accept")
	}
	return visitor.VisitExpressionStmt(e)
}

func NewExpressionStmt(expr Expr) *Expression {
	return &Expression{Expr: expr}
}

// Print type
type Print struct {
	Expr Expr
}

func (p *Print) Accept(visitor StmtVisitor) interface{} {
	if visitor == nil {
		panic("Visitor is nil in Print.Accept")
	}
	return visitor.VisitPrintStmt(p)
}

func NewPrintStmt(expr Expr) *Print {
	return &Print{Expr: expr}
}
