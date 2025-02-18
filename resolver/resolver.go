package resolver

import (
	"github.com/drewslam/goloxTreeInterpreter/ast"
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

func (r *Resolver) resolve(input interface{}) {
	switch v := input.(type) {
	case ast.Stmt:
		v.Accept(r)
	case ast.Expr:
		v.Accept(r)
	default:
	}
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

func (r *Resolver) VisitVarStmt(stmt *ast.Var) interface{} {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolve(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}
