package parser

import (
	"fmt"

	"github.com/drewslam/goloxTreeInterpreter/ast"
	"github.com/drewslam/goloxTreeInterpreter/errors"
	"github.com/drewslam/goloxTreeInterpreter/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func (p *Parser) NewParser(tokens []token.Token) {
	p.tokens = tokens
	p.current = 0
}

func (p *Parser) expression() ast.Expr {
	return p.assignment()
}

func (p *Parser) declaration() ast.Stmt {
	defer func() {
		if r := recover(); r != nil {
			// Check if it's a parse error
			if _, ok := r.(*errors.ParseError); ok {
				// Return nil if a ParseError is caught
				return
			}
			// Re-panic for unexpected errors
			panic(r)
		}
	}()

	if p.match(token.CLASS) {
		return p.classDeclaration()
	}
	if p.match(token.FUN) {
		return p.function("function")
	}
	if p.match(token.VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) classDeclaration() ast.Stmt {
	name := p.consume(token.IDENTIFIER, "Expect class name.")
	p.consume(token.LEFT_BRACE, "Expect '{' before class body.")

	var methods []ast.Function
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		methods = append(methods, *p.function("method"))
	}

	p.consume(token.RIGHT_BRACE, "Expect '}' after class body.")
	return &ast.Class{
		Name:    name,
		Methods: methods,
	}
}

func (p *Parser) statement() ast.Stmt {
	if p.match(token.FOR) {
		return p.forStatement()
	}
	if p.match(token.IF) {
		return p.ifStatement()
	}
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	if p.match(token.RETURN) {
		return p.returnStatement()
	}
	if p.match(token.WHILE) {
		return p.whileStatement()
	}
	if p.match(token.LEFT_BRACE) {
		return &ast.Block{Statements: p.block()}
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() ast.Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer ast.Stmt
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition ast.Expr = nil
	if !p.check(token.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(token.SEMICOLON, "Expect ';' after loop condition.")

	var increment ast.Expr = nil
	if !p.check(token.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")

	body := p.statement()

	if increment != nil {
		body = &ast.Block{
			Statements: []ast.Stmt{
				body,
				&ast.Expression{Expr: increment},
			},
		}
	}

	/*	if condition == nil {
			condition = &ast.Literal{
				Value: true,
			}
		}
	*/
	body = &ast.While{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = &ast.Block{
			Statements: []ast.Stmt{
				initializer,
				body,
			},
		}
	}

	return body
}

func (p *Parser) ifStatement() ast.Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()

	var elseBranch ast.Stmt = nil
	if p.match(token.ELSE) {
		elseBranch = p.statement()
	}

	return ast.NewIfStmt(condition, thenBranch, elseBranch)
}

func (p *Parser) printStatement() ast.Stmt {
	expr := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return ast.NewPrintStmt(expr)
}

func (p *Parser) returnStatement() ast.Stmt {
	keyword := p.previous()
	var value ast.Expr = nil
	if !p.check(token.SEMICOLON) {
		value = p.expression()
	}

	p.consume(token.SEMICOLON, "Expect ';' after return value.")
	return &ast.Return{
		Keyword: keyword,
		Value:   value,
	}
}

func (p *Parser) varDeclaration() ast.Stmt {
	name := p.consume(token.IDENTIFIER, "Expect variable name.")

	var initializer ast.Expr
	if p.match(token.EQUAL) {
		initializer = p.expression()
	}

	p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
	return ast.NewVarStmt(name, initializer)
}

func (p *Parser) whileStatement() ast.Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "Expect ')' after condition.")

	body := p.statement()

	/*	body = &ast.Block{
			Statements: []ast.Stmt{
				body,
			},
		}
	*/
	return &ast.While{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) expressionStatement() ast.Stmt {
	expr := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after expression.")
	return ast.NewExpressionStmt(expr)
}

func (p *Parser) function(kind string) *ast.Function {
	message := fmt.Sprintf("Expect %v name.", kind)
	name := p.consume(token.IDENTIFIER, message)

	p.consume(token.LEFT_PAREN, "Expect '(' after function name.")

	var parameters []token.Token
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				errors.ReportParseError(p.peek(), "Can't have more than 255 parameters.")
			}
			parameters = append(parameters, p.consume(token.IDENTIFIER, "Expect parameter name."))
			if !p.match(token.COMMA) {
				break
			}
		}
	}

	p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")

	message = fmt.Sprintf("Expect '{' before %v body.", kind)
	p.consume(token.LEFT_BRACE, message)

	body := p.block()
	return &ast.Function{
		Name:   name,
		Params: parameters,
		Body:   body,
	}
}

func (p *Parser) block() []ast.Stmt {
	var statements []ast.Stmt

	// p.consume(token.LEFT_BRACE, "Expect '{' before block.")

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(token.RIGHT_BRACE, "Expect '}' after block.")

	return statements
}

func (p *Parser) assignment() ast.Expr {
	expr := p.or()

	if p.match(token.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if variable, ok := expr.(*ast.Variable); ok {
			name := variable.Name
			return &ast.Assign{
				Name:  name,
				Value: value,
			}
		}
		if get, ok := expr.(*ast.Get); ok {
			return &ast.Set{
				Object: get.Object,
				Name:   get.Name,
				Value:  value,
			}
		}

		errors.ReportParseError(equals, "Invalid assignment target.")
	}

	return expr
}

func (p *Parser) or() ast.Expr {
	expr := p.and()

	for p.match(token.OR) {
		operator := p.previous()
		right := p.and()
		expr = &ast.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) and() ast.Expr {
	expr := p.equality()

	for p.match(token.AND) {
		operator := p.previous()
		right := p.equality()
		expr = &ast.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) equality() ast.Expr {
	expr := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) comparison() ast.Expr {
	expr := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) unary() ast.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return &ast.Unary{
			Operator: operator,
			Right:    right,
		}
	}

	return p.call()
}

func (p *Parser) finishCall(callee ast.Expr) ast.Expr {
	var arguments []ast.Expr

	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				errors.ReportParseError(p.peek(), "Can't have more than 255 arguments.")
			}

			arguments = append(arguments, p.expression())

			if !p.match(token.COMMA) {
				break
			}
		}
	}

	paren := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")

	return &ast.Call{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}
}

func (p *Parser) call() ast.Expr {
	expr := p.primary()

	for {
		if p.match(token.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(token.DOT) {
			name := p.consume(token.IDENTIFIER, "Expect property name after '.'.")
			expr = &ast.Get{
				Object: expr,
				Name:   name,
			}
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) primary() ast.Expr {
	if p.match(token.FALSE) {
		return &ast.Literal{
			Value: false,
		}
	}
	if p.match(token.TRUE) {
		return &ast.Literal{
			Value: true,
		}
	}
	if p.match(token.NIL) {
		return &ast.Literal{
			Value: nil,
		}
	}

	if p.match(token.NUMBER, token.STRING) {
		return &ast.Literal{
			Value: p.previous().Literal,
		}
	}

	if p.match(token.IDENTIFIER) {
		return &ast.Variable{
			Name: p.previous(),
		}
	}

	if p.match(token.LEFT_PAREN) {
		expr := p.expression()

		p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		return &ast.Grouping{
			Expression: expr,
		}
	}

	// Error handling if no valid expression is found
	errors.ReportParseError(p.peek(), "Expect expression.")
	return nil
}

func (p *Parser) Parse() []ast.Stmt {
	var statements []ast.Stmt

	defer func() {
		if r := recover(); r != nil {
			// Check if it's a parse error
			if _, ok := r.(*errors.ParseError); ok {
				// Return nil if a ParseError is caught
				return
			}
			// Re-panic for unexpected errors
			panic(r)
		}
	}()

	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokentype token.TokenType, message string) token.Token {
	if p.check(tokentype) {
		return p.advance()
	}

	errors.ReportError(p.peek().Line, message)
	p.synchronize()

	return token.Token{}
}

func (p *Parser) check(tokentype token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokentype
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}
