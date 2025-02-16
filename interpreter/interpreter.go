package interpreter

import (
	"fmt"

	"github.com/drewslam/goloxTreeInterpreter/ast"
	"github.com/drewslam/goloxTreeInterpreter/environment"
	"github.com/drewslam/goloxTreeInterpreter/errors"
	"github.com/drewslam/goloxTreeInterpreter/token"
)

type Interpreter struct {
	object      ast.ExprVisitor
	void        ast.StmtVisitor
	environment *environment.Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: environment.NewEnvironment(),
	}
}

func (i *Interpreter) Interpret(statements []ast.Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if runtimeErr, ok := r.(*errors.RuntimeError); ok {
				i.reportRuntimeError(runtimeErr)
			} else {
				panic(r) // Re-panic if it's not a RuntimeError
			}
		}
	}()

	// value := i.evaluate(expr)
	// fmt.Println(i.stringify(value))
	for _, stmt := range statements {
		i.execute(stmt)
	}
}

func (i *Interpreter) execute(stmt ast.Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, environment *environment.Environment) {
	previous := i.environment
	i.environment = environment
	for _, statement := range statements {
		i.execute(statement)
	}
	i.environment = previous
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.Block) interface{} {
	i.executeBlock(stmt.Statements, environment.NewEnvironment(i.environment))
	return nil
}

func (i *Interpreter) evaluate(expr ast.Expr) interface{} {
	if expr == nil {
		panic("Tried to evaluate a nil expression.")
	}
	return expr.Accept(i)
}

func (i *Interpreter) VisitExpressionStmt(stmt *ast.Expression) interface{} {
	i.evaluate(stmt.Expr)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *ast.If) interface{} {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.Print) interface{} {
	value := i.evaluate(stmt.Expr)
	fmt.Println(i.stringify(value))
	return nil
}

func (i *Interpreter) VisitVarStmt(stmt *ast.Var) interface{} {
	var value interface{} = nil
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}

	i.environment.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.While) interface{} {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
	return nil
}

func (i *Interpreter) VisitAssignExpr(expr *ast.Assign) interface{} {
	value := i.evaluate(expr.Value)
	if err := i.environment.Assign(expr.Name, value); err != nil {
		panic(errors.NewRuntimeError(expr.Name, "Undefined variable '"+expr.Name.Lexeme+"'"))
	}
	return value
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.BANG_EQUAL:
		return !i.isEqual(left, right)
	case token.EQUAL_EQUAL:
		return i.isEqual(left, right)
	case token.GREATER:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case token.GREATER_EQUAL:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case token.LESS:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case token.LESS_EQUAL:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case token.MINUS:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)
	case token.PLUS:
		if leftVal, ok := left.(float64); ok {
			if rightVal, ok := right.(float64); ok {
				return leftVal + rightVal
			}
		}
		if leftVal, ok := left.(string); ok {
			if rightVal, ok := right.(string); ok {
				return leftVal + rightVal
			}
		}
		panic(errors.NewRuntimeError(expr.Operator, "Operands must be two numbers or two strings."))
	case token.SLASH:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) / right.(float64)
	case token.STAR:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) * right.(float64)
	}

	// Unreachable
	return nil
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping) interface{} {
	// Handle grouping expressions
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.Literal) interface{} {
	// Handle literal expressions
	return expr.Value
}

func (i *Interpreter) VisitLogicalExpr(expr *ast.Logical) interface{} {
	// Handle logical expressions
	left := i.evaluate(expr.Left)

	if expr.Operator.Type == token.OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) interface{} {
	// Handle unary expressions
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.BANG:
		return !i.isTruthy(right)
	case token.MINUS:
		i.checkNumberOperand(expr.Operator, right)
		return -right.(float64)
	}

	// Unreachable
	return nil
}

func (i *Interpreter) VisitVariableExpr(expr *ast.Variable) interface{} {
	return i.environment.Get(expr.Name)
}

func (i *Interpreter) checkNumberOperand(operator token.Token, operand interface{}) {
	if _, ok := operand.(float64); ok {
		return
	}
	panic(errors.NewRuntimeError(operator, "Operand must be a number."))
}

func (i *Interpreter) checkNumberOperands(operator token.Token, left interface{}, right interface{}) {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return
		}
	}
	panic(errors.NewRuntimeError(operator, "Operands must be two numbers."))
}

func (i *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if value, ok := object.(bool); ok {
		return value
	}
	return true
}

func (i *Interpreter) isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return a == b
}

// stringify converts an evaluated object into a human-readable string
func (i *Interpreter) stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}

	if val, ok := object.(float64); ok {
		text := fmt.Sprintf("%g", val)
		return text
	}

	return fmt.Sprintf("%v", object)
}

// reportRuntimeError handles runtime error reporting
func (i *Interpreter) reportRuntimeError(err *errors.RuntimeError) {
	fmt.Printf("[line %d] RuntimeError: %s\n", err.Token.Line, err.Message)
	errors.HadRuntimeError = true
}
