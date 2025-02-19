package resolver

import (
	"github.com/drewslam/goloxTreeInterpreter/ast"
	"github.com/drewslam/goloxTreeInterpreter/errors"
	"github.com/drewslam/goloxTreeInterpreter/interpreter"
	"github.com/drewslam/goloxTreeInterpreter/token"
)

func peek(stack []map[string]bool) (map[string]bool, bool) {
	if len(stack) == 0 {
		return nil, false
	}
	return stack[len(stack)-1], true
}

type Resolver struct {
	interpreter *interpreter.Interpreter
	scopes      []map[string]bool
}

var _ ast.StmtVisitor = (*Resolver)(nil)
var _ ast.ExprVisitor = (*Resolver)(nil)

func (r *Resolver) Resolve(statements []ast.Stmt) {
	for _, statement := range statements {
		r.resolve(statement)
	}
}

func (r *Resolver) VisitBlockStmt(stmt *ast.Block) interface{} {
	r.beginScope()
	r.Resolve(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ast.Expression) interface{} {
	r.resolve(stmt.Expr)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *ast.Function) interface{} {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *ast.If) interface{} {
	r.resolve(stmt.Condition)
	r.resolve(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolve(stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *ast.Print) interface{} {
	r.resolve(stmt.Expr)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *ast.Return) interface{} {
	if stmt.Value != nil {
		r.resolve(stmt.Value)
	}
	return nil
}

func (r *Resolver) resolve(input interface{}) {
	switch v := input.(type) {
	case ast.Stmt:
		v.Accept(r)
	case ast.Expr:
		v.Accept(r)
	default:
	}
}

func (r *Resolver) resolveFunction(function *ast.Function) {
	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	r.resolve(function.Body)
	r.endScope()
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope, ok := peek(r.scopes)
	if ok {
		scope[name.Lexeme] = false
	}
}

func (r *Resolver) define(name token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope, ok := peek(r.scopes)
	if ok {
		scope[name.Lexeme] = true
	}
}

func (r *Resolver) resolveLocal(expr ast.Expr, name token.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.Resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) VisitVarStmt(stmt *ast.Var) interface{} {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolve(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *ast.While) interface{} {
	r.resolve(stmt.Condition)
	r.resolve(stmt.Body)
	return nil
}

func (r *Resolver) VisitAssignExpr(expr *ast.Assign) interface{} {
	r.resolve(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr *ast.Binary) interface{} {
	r.resolve(expr.Left)
	r.resolve(expr.Right)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *ast.Call) interface{} {
	r.resolve(expr.Callee)
	for _, argument := range expr.Arguments {
		r.resolve(argument)
	}
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *ast.Grouping) interface{} {
	r.resolve(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *ast.Literal) interface{} {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr *ast.Logical) interface{} {
	r.resolve(expr.Left)
	r.resolve(expr.Right)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *ast.Unary) interface{} {
	r.resolve(expr.Right)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *ast.Variable) interface{} {
	if len(r.scopes) > 0 {
		if scope, ok := peek(r.scopes); ok {
			if defined, exists := scope[expr.Name.Lexeme]; exists && !defined {
				errors.ReportParseError(expr.Name, "Can't read local variable in its own initializer.")
			}
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil
}
