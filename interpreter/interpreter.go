package interpreter

import (
	"fmt"

	"github.com/drewslam/goloxTreeInterpreter/ast"
	"github.com/drewslam/goloxTreeInterpreter/environment"
	"github.com/drewslam/goloxTreeInterpreter/loxCallable"
	"github.com/drewslam/goloxTreeInterpreter/loxError"
	"github.com/drewslam/goloxTreeInterpreter/object"
	"github.com/drewslam/goloxTreeInterpreter/returnValue"
	"github.com/drewslam/goloxTreeInterpreter/token"
)

type Interpreter struct {
	exprVisitor ast.ExprVisitor
	stmtVisitor ast.StmtVisitor
	Globals     *environment.Environment
	locals      map[ast.Expr]int
	environment *environment.Environment
}

func NewInterpreter() *Interpreter {
	globalEnv := environment.NewEnvironment()

	loxCallable.RegisterNatives(globalEnv)

	return &Interpreter{
		Globals:     globalEnv,
		environment: globalEnv,
		locals:      make(map[ast.Expr]int),
	}
}

func (i *Interpreter) StoreResolution(expr ast.Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) Interpret(statements []ast.Stmt) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case *loxError.LoxError:
				if v.IsFatal {
					loxError.ReportAndPanic(v)
				}
			case *returnValue.ReturnValue:
				if v.Value != nil {
					fmt.Println(v.Value)
				}
			default:
				panic(r) // Re-panic if it's not a RuntimeError
			}
		}
	}()

	for _, stmt := range statements {
		/*if err := i.execute(stmt); err != nil {
			panic(err)
		}*/
		result := i.execute(stmt)
		fmt.Printf("Environment after executing %T: %+v\n", stmt, i.environment.Values)
		if result != nil {
			fmt.Printf("Executiong returned unexpected value: %v\n", result)
		}
	}
}

func (i *Interpreter) execute(stmt ast.Stmt) interface{} {
	result := stmt.Accept(i)
	fmt.Printf("Executing: %T -> result: %v\n", stmt, result)
	return result
}

func (i *Interpreter) Resolve(expr ast.Expr, depth int) {
	if i.locals == nil {
		panic("Interpreter.locals is nil!")
	}
	i.locals[expr] = depth
}

func (i *Interpreter) GetGlobals() *environment.Environment {
	return i.Globals
}

func (i *Interpreter) ExecuteBlock(statements []ast.Stmt, environment *environment.Environment) interface{} {
	previous := i.environment
	i.environment = environment

	defer func() { i.environment = previous }()

	defer func() {
		if r := recover(); r != nil {
			if returnValue, ok := r.(*returnValue.ReturnValue); ok {
				i.environment = previous
				panic(returnValue)
			}
			panic(r)
		}
	}()

	for _, statement := range statements {
		result := i.execute(statement)

		if returnVal, ok := result.(*returnValue.ReturnValue); ok {
			return returnVal.Value
		}
	}
	return nil
}

var _ loxCallable.Interpreter = &Interpreter{}

func (i *Interpreter) VisitBlockStmt(stmt *ast.Block) interface{} {
	err := i.ExecuteBlock(stmt.Statements, environment.NewEnvironment(i.environment))
	if err != nil {
		return err
	}
	return nil
}

func (i *Interpreter) VisitClassStmt(stmt *ast.Class) interface{} {
	i.environment.Define(stmt.Name.Lexeme, nil)

	methods := make(map[string]*object.LoxFunction)
	for _, method := range stmt.Methods {
		function := object.NewLoxFunction(method, i.environment, method.Name.Lexeme == "init")
		methods[method.Name.Lexeme] = function
	}

	klass := &object.LoxClass{
		Name:    stmt.Name.Lexeme,
		Methods: methods,
	}
	/*
		instance := &object.LoxInstance{
			Klass:  klass,
			Fields: make(map[string]interface{}),
		}
	*/
	i.environment.Assign(stmt.Name, klass)
	return nil
}

func (i *Interpreter) evaluate(expr ast.Expr) interface{} {
	if expr == nil {
		err := loxError.NewRuntimeError(token.Token{Line: 0}, "", "Tried to evaluate a nil expression.")
		loxError.ReportAndPanic(err)
	}

	fmt.Printf("Evaluating expression: %T\n", expr)

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(*loxError.LoxError); ok && err.IsFatal {
				loxError.ReportAndPanic(err)
			}
		}
	}()

	result := expr.Accept(i)
	fmt.Printf("Expression result: %v (type: %T)\n", result, result)
	return result
}

func (i *Interpreter) VisitExpressionStmt(stmt *ast.Expression) interface{} {
	i.evaluate(stmt.Expr)
	return nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *ast.Function) interface{} {
	function := object.NewLoxFunction(stmt, i.environment, false)
	i.environment.Define(stmt.Name.Lexeme, function)
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
	fmt.Printf("Printing value: %v\n", value)
	fmt.Println(i.stringify(value))
	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.Return) interface{} {
	var value interface{} = nil
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}
	fmt.Println("Panic with return value:", value)
	panic(&returnValue.ReturnValue{Value: value})
}

func (i *Interpreter) VisitVarStmt(stmt *ast.Var) interface{} {
	var value interface{} = nil
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}

	fmt.Printf("Storing variable '%s' with value: %v (type: %T)\n", stmt.Name.Lexeme, value, value)
	i.environment.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.While) interface{} {
	previous := i.environment

	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}

	i.environment = previous
	return nil
}

func (i *Interpreter) VisitAssignExpr(expr *ast.Assign) interface{} {
	value := i.evaluate(expr.Value)

	if distance, exists := i.locals[expr]; exists {
		i.environment.AssignAt(distance, expr.Name, value)
	} else {
		i.Globals.Assign(expr.Name, value)
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
		err := loxError.NewRuntimeError(expr.Operator, fmt.Sprintf("[Line %d]: ", expr.Operator.Line), "Operands must be two numbers or two strings.")
		loxError.ReportAndPanic(err)
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

func (i *Interpreter) VisitCallExpr(expr *ast.Call) interface{} {
	callee := i.evaluate(expr.Callee)

	fmt.Printf("Calling function: %v (type: %T)\n", callee, callee)

	var arguments []interface{}
	for _, argument := range expr.Arguments {
		evaluatedArg := i.evaluate(argument)
		fmt.Printf("Evaluated arguments: %v (type: %T)\n", evaluatedArg, evaluatedArg)
		arguments = append(arguments, evaluatedArg)
	}

	function, ok := callee.(loxCallable.LoxCallable)
	if !ok {
		err := loxError.NewRuntimeError(expr.Paren, fmt.Sprintf("[Line %d]: ", expr.Paren.Line), "Can only call functions and classes.")
		loxError.ReportAndPanic(err)
	}

	if len(arguments) != function.Arity() {
		message := fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments))
		err := loxError.NewRuntimeError(expr.Paren, fmt.Sprintf("[Line %d]: ", expr.Paren.Line), message)
		loxError.ReportAndPanic(err)
	}

	result := function.Call(i, arguments)
	fmt.Printf("Function returned: %v (type: %T)\n", result, result)
	return result
}

func (i *Interpreter) VisitGetExpr(expr *ast.Get) interface{} {
	objekt := i.evaluate(expr.Object)
	if instance, ok := objekt.(*object.LoxInstance); ok {
		return instance.Get(expr.Name)
	}

	return loxError.NewRuntimeError(expr.Name, fmt.Sprintf("[Line %d]: ", expr.Name.Line), "Only instances have properties.")
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

func (i *Interpreter) VisitSetExpr(expr *ast.Set) interface{} {
	objekt := i.evaluate(expr.Object)

	if _, ok := objekt.(*object.LoxInstance); !ok {
		return loxError.NewRuntimeError(expr.Name, fmt.Sprintf("[Line %d]: ", expr.Name.Line), "Only instances have fields.")
	}

	value := i.evaluate(expr.Value)
	objekt.(*object.LoxInstance).Set(expr.Name, value)
	return value
}

func (i *Interpreter) VisitThisExpr(expr *ast.This) interface{} {
	return i.lookUpVariable(expr.Keyword, expr)
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
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) lookUpVariable(name token.Token, expr ast.Expr) interface{} {
	if i.environment != nil {
		val, err := i.environment.Get(name)
		if err == nil {
			return val
		}
	}

	distance, exists := i.locals[expr]

	if !exists {
		for lookupExpr, lookupDistance := range i.locals {
			if varExpr, ok := lookupExpr.(*ast.Variable); ok &&
				varExpr.Name.Lexeme == name.Lexeme {
				distance = lookupDistance
				exists = true
				break
			}
		}
	}

	if exists {
		fmt.Printf("Looking up local variable '%s' at distance %d\n", name.Lexeme, distance)
		res, err := i.environment.GetAt(distance, name.Lexeme)
		if err != nil {
			fmt.Printf("Error retrieving local variable '%s': %v\n", name.Lexeme, err)
			err := loxError.NewRuntimeError(name, fmt.Sprintf("%d", name.Line), "Undefined local variable '"+name.Lexeme+"'")
			loxError.ReportAndPanic(err)
		}
		fmt.Printf("Local variable '%s' resolved to: %v\n", name.Lexeme, res)
		return res
	}

	// Fallback to globals
	fmt.Printf("Variable '%s' not found locally. Checking globals\n", name.Lexeme)
	res, err := i.Globals.Get(name)
	if err != nil {
		fmt.Printf("Error retrieving global variable '%s': %v\n", name.Lexeme, err)
		err := loxError.NewRuntimeError(name, fmt.Sprintf("%d", name.Line), "Undefined ariable '"+name.Lexeme+"'")
		loxError.ReportAndPanic(err)
	}
	fmt.Printf("Global variable '%s' resolved to: %v\n", name.Lexeme, res)
	return res
}

func (i *Interpreter) checkNumberOperand(operator token.Token, operand interface{}) {
	if _, ok := operand.(float64); ok {
		return
	}
	panic(loxError.NewRuntimeError(operator, fmt.Sprintf("[Line %d]: ", operator.Line), "Operand must be a number."))
}

func (i *Interpreter) checkNumberOperands(operator token.Token, left interface{}, right interface{}) {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return
		}
	}
	panic(loxError.NewRuntimeError(operator, fmt.Sprintf("[Line %d]: ", operator.Line), "Operands must be two numbers."))
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

/*
// reportRuntimeError handles runtime error reporting
func (i *Interpreter) reportRuntimeError(err *loxError.RuntimeError) {
	fmt.Printf("[line %d] RuntimeError: %s\n", err.Token.Line, err.Message)
	loxError.HadRuntimeError = true
}*/
