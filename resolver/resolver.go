package resolver

import (
	"fmt"

	"github.com/drewslam/goloxTreeInterpreter/ast"
	"github.com/drewslam/goloxTreeInterpreter/interpreter"
	"github.com/drewslam/goloxTreeInterpreter/loxDebug"
	"github.com/drewslam/goloxTreeInterpreter/loxError"
	"github.com/drewslam/goloxTreeInterpreter/token"
)

func peek(stack []map[string]bool) (map[string]bool, bool) {
	if len(stack) == 0 {
		return nil, false
	}
	return stack[len(stack)-1], true
}

type Resolver struct {
	Interpreter     *interpreter.Interpreter
	scopes          []map[string]bool
	CurrentFunction FunctionType
}

type FunctionType int

const (
	NOT_FUNCTION FunctionType = iota
	FUNCTION
	INITIALIZER
	METHOD
)

func NewResolver(interpreter *interpreter.Interpreter) *Resolver {
	return &Resolver{
		Interpreter:     interpreter,
		scopes:          make([]map[string]bool, 0),
		CurrentFunction: NOT_FUNCTION,
	}
}

type ClassType int

const (
	NOT_CLASS ClassType = iota
	CLASS
)

var currentClass ClassType = NOT_CLASS

var _ ast.StmtVisitor = (*Resolver)(nil)
var _ ast.ExprVisitor = (*Resolver)(nil)

func (r *Resolver) Resolve(statements []ast.Stmt) *loxError.LoxError {
	for _, statement := range statements {
		r.resolve(statement)
	}
	return nil
}

func (r *Resolver) VisitBlockStmt(stmt *ast.Block) interface{} {
	r.beginScope()
	r.Resolve(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitClassStmt(stmt *ast.Class) interface{} {
	enclosingClass := currentClass
	currentClass = CLASS

	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.beginScope()
	if len(r.scopes) > 0 {
		r.scopes[len(r.scopes)-1]["this"] = true
	}

	for _, method := range stmt.Methods {
		declaration := METHOD
		if method.Name.Lexeme == "init" {
			declaration = INITIALIZER
		}
		r.resolveFunction(method, declaration)
	}
	r.endScope()

	currentClass = enclosingClass
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ast.Expression) interface{} {
	r.resolve(stmt.Expr)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *ast.Function) interface{} {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, FUNCTION)
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
	if r.CurrentFunction == NOT_FUNCTION {
		return loxError.NewRuntimeError(stmt.Keyword, fmt.Sprintf("[Line %d]: ", stmt.Keyword.Line), "Can't return from top level code.'")
	}

	if stmt.Value != nil {
		if r.CurrentFunction == INITIALIZER {
			return loxError.NewRuntimeError(stmt.Keyword, fmt.Sprintf("[Line %d]: ", stmt.Keyword.Line), "Can't return a value from an initializer.")
		}

		r.resolve(stmt.Value)
	}
	return nil
}

func (r *Resolver) resolve(input interface{}) {
	if len(r.scopes) == 0 {
		r.beginScope()
	}

	switch v := input.(type) {
	case ast.Stmt:
		v.Accept(r)
	case ast.Expr:
		v.Accept(r)
	default:
	}
}

func (r *Resolver) resolveFunction(function *ast.Function, functiontype FunctionType) {
	enclosingFunction := r.CurrentFunction
	r.CurrentFunction = functiontype

	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	r.resolve(function.Body)
	r.endScope()

	r.CurrentFunction = enclosingFunction
}

func (r *Resolver) beginScope() {
	if r.scopes == nil {
		r.scopes = []map[string]bool{}
	}

	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	if len(r.scopes) > 0 {
		r.scopes = r.scopes[:len(r.scopes)-1]
	}
}

func (r *Resolver) declare(name token.Token) *loxError.LoxError {
	if len(r.scopes) == 0 {
		return nil
	}

	if scope, ok := peek(r.scopes); ok {
		scope[name.Lexeme] = false
	} else {
		return loxError.NewRuntimeError(name, fmt.Sprintf("[Line %d]: ", name.Line), "Already a variable with this name in this scope.")
	}

	return nil
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
	/*if len(r.scopes) == 0 {
		// panic("No active scope when resolving a variable.")
		return
	}*/
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			depth := len(r.scopes) - 1 - i
			if i == 0 {
				loxDebug.LogInfo("Variable '%s' correctly marked as global\n", name.Lexeme)
			}
			loxDebug.LogInfo("Resolving variable '%s' as local at depth %d\n", name.Lexeme, depth)
			r.Interpreter.Resolve(expr, depth)
			// r.Interpreter.StoreResolution(expr, depth)
			return
		}
	}
	loxDebug.LogInfo("Variable '%s' is treated as global\n", name.Lexeme)
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

	if expr.Name.Lexeme == "this" {
		err := loxError.NewScanError(expr.Name.Line, "Invalid assignment target.")
		loxError.ReportError(err)
	}

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

func (r *Resolver) VisitGetExpr(expr *ast.Get) interface{} {
	r.resolve(expr.Object)
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

func (r *Resolver) VisitSetExpr(expr *ast.Set) interface{} {
	r.resolve(expr.Value)
	r.resolve(expr.Object)
	return nil
}

func (r *Resolver) VisitThisExpr(expr *ast.This) interface{} {
	if currentClass == NOT_CLASS {
		return loxError.NewRuntimeError(expr.Keyword, fmt.Sprintf("[Line %d]: ", expr.Keyword.Line), "Can't use 'this' outside of a class.")
	}

	r.resolveLocal(expr, expr.Keyword)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *ast.Unary) interface{} {
	r.resolve(expr.Right)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *ast.Variable) interface{} {
	loxDebug.LogDebug("Current scopes:", r.scopes)

	if len(r.scopes) > 0 {
		if scope, ok := peek(r.scopes); ok {
			if defined, exists := scope[expr.Name.Lexeme]; exists {
				if !defined {
					return loxError.NewParseError(expr.Name, "Can't read local variable in its own initializer.")
				}
			}
		}
	} else {
		panic("No active scope when resolving a variable.")
	}

	r.resolveLocal(expr, expr.Name)
	return nil
}
