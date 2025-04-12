package parser

import (
	"fmt"

	"github.com/drewslam/goloxTreeInterpreter/ast"
	"github.com/drewslam/goloxTreeInterpreter/loxError"
	"github.com/drewslam/goloxTreeInterpreter/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) expression() (ast.Expr, *loxError.LoxError) {
	return p.assignment()
}

func (p *Parser) declaration() (ast.Stmt, *loxError.LoxError) {
	defer func() {
		if r := recover(); r != nil {
			// Check if it's a parse error
			if _, ok := r.(*loxError.LoxError); ok {
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

func (p *Parser) classDeclaration() (ast.Stmt, *loxError.LoxError) {
	name := p.consume(token.IDENTIFIER, "Expect class name.")

	var superclass *ast.Variable
	if p.match(token.LESS) {
		p.consume(token.IDENTIFIER, "Expect superclass name.")
		superclass = &ast.Variable{
			Name: p.previous(),
		}
	}

	p.consume(token.LEFT_BRACE, "Expect '{' before class body.")

	var methods []*ast.Function
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		if !p.check(token.IDENTIFIER) {
			// err := loxError.NewParseError(p.peek(), "Only methods are allowed in class bodies.")
			p.synchronize()
			continue
		}

		method, err := p.function("method")
		if err != nil {
			return nil, err
		} else {
			methods = append(methods, method)
		}

	}

	p.consume(token.RIGHT_BRACE, "Expect '}' after class body.")
	return &ast.Class{
		Name:       name,
		Superclass: superclass,
		Methods:    methods,
	}, nil
}

func (p *Parser) statement() (ast.Stmt, *loxError.LoxError) {
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
		block, err := p.block()
		if err != nil {
			return nil, err
		}
		return ast.NewBlockStmt(block), nil
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() (ast.Stmt, *loxError.LoxError) {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer ast.Stmt
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		val, err := p.varDeclaration()
		if err != nil {
			return nil, err
		}
		initializer = val
	} else {
		val, err := p.expressionStatement()
		if err != nil {
			return nil, err
		}
		initializer = val
	}

	var condition ast.Expr = nil
	if !p.check(token.SEMICOLON) {
		val, err := p.expression()
		if err != nil {
			return nil, err
		}
		condition = val
	}
	p.consume(token.SEMICOLON, "Expect ';' after loop condition.")

	var increment ast.Expr = nil
	if !p.check(token.RIGHT_PAREN) {
		val, err := p.expression()
		if err != nil {
			return nil, err
		}
		increment = val
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = &ast.Block{
			Statements: []ast.Stmt{
				body,
				&ast.Expression{Expr: increment},
			},
		}
	}

	if condition == nil {
		condition = &ast.Literal{
			Value: true,
		}
	}

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

	return body, nil
}

func (p *Parser) ifStatement() (ast.Stmt, *loxError.LoxError) {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch ast.Stmt = nil
	if p.match(token.ELSE) {
		val, err := p.statement()
		if err != nil {
			return nil, err
		}
		elseBranch = val
	}

	return ast.NewIfStmt(condition, thenBranch, elseBranch), nil
}

func (p *Parser) printStatement() (ast.Stmt, *loxError.LoxError) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return ast.NewPrintStmt(expr), nil
}

func (p *Parser) returnStatement() (ast.Stmt, *loxError.LoxError) {
	keyword := p.previous()
	var value ast.Expr = nil
	if !p.check(token.SEMICOLON) {
		val, err := p.expression()
		if err != nil {
			return nil, err
		}
		value = val
	}

	p.consume(token.SEMICOLON, "Expect ';' after return value.")
	return &ast.Return{
		Keyword: keyword,
		Value:   value,
	}, nil
}

func (p *Parser) varDeclaration() (ast.Stmt, *loxError.LoxError) {
	name := p.consume(token.IDENTIFIER, "Expect variable name.")

	var initializer ast.Expr
	if p.match(token.EQUAL) {
		val, err := p.expression()
		if err != nil {
			return nil, loxError.NewParseError(name, "Invalid variable declaration.")
		}
		initializer = val
	}

	p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
	return ast.NewVarStmt(name, initializer), nil
}

func (p *Parser) whileStatement() (ast.Stmt, *loxError.LoxError) {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after condition.")

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &ast.While{
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) expressionStatement() (ast.Stmt, *loxError.LoxError) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(token.SEMICOLON, "Expect ';' after expression.")
	return ast.NewExpressionStmt(expr), nil
}

func (p *Parser) function(kind string) (*ast.Function, *loxError.LoxError) {
	message := fmt.Sprintf("Expect %v name.", kind)
	name := p.consume(token.IDENTIFIER, message)

	p.consume(token.LEFT_PAREN, "Expect '(' after function name.")

	var parameters []token.Token
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				return nil, loxError.NewParseError(p.peek(), "Can't have more than 255 parameters.")
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

	body, err := p.block()
	if err != nil {
		return nil, err
	}
	return &ast.Function{
		Name:   name,
		Params: parameters,
		Body:   body,
	}, nil
}

func (p *Parser) block() ([]ast.Stmt, *loxError.LoxError) {
	var statements []ast.Stmt

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		declaration, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, declaration)
	}

	p.consume(token.RIGHT_BRACE, "Expect '}' after block.")

	return statements, nil
}

func (p *Parser) assignment() (ast.Expr, *loxError.LoxError) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}
	/*
		if expr == nil {
			loxDebug.LogError("Parsed a nil expression in assignment().")
		}
	*/
	if p.match(token.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()

		if err != nil {
			// p.synchronize()
			return nil, err
		}

		switch v := expr.(type) {
		case *ast.Variable:
			name := v.Name
			if name.Lexeme == "this" {
				// p.synchronize()
				return nil, loxError.NewParseError(name, "Cannot assign to 'this'.")
			}
			return &ast.Assign{
				Name:  name,
				Value: value,
			}, nil
		case *ast.Get:
			return &ast.Set{
				Object: v.Object,
				Name:   v.Name,
				Value:  value,
			}, nil
		case *ast.This:
			return nil, loxError.NewParseError(v.Keyword, "Cannot assign to 'this'.")
		default:
			return nil, loxError.NewParseError(equals, "Invalid assignment target.")
		}

		// p.synchronize()
		// return nil, loxError.NewParseError(equals, "Invalid assignment target.")
	}

	return expr, nil
}

func (p *Parser) or() (ast.Expr, *loxError.LoxError) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(token.OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &ast.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) and() (ast.Expr, *loxError.LoxError) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(token.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &ast.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) equality() (ast.Expr, *loxError.LoxError) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, *loxError.LoxError) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (ast.Expr, *loxError.LoxError) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) factor() (ast.Expr, *loxError.LoxError) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) unary() (ast.Expr, *loxError.LoxError) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ast.Unary{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.call()
}

func (p *Parser) finishCall(callee ast.Expr) (ast.Expr, *loxError.LoxError) {
	var arguments []ast.Expr

	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				return nil, loxError.NewParseError(p.peek(), "Can't have more than 255 arguments.")
			}

			argument, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, argument)

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
	}, nil
}

func (p *Parser) call() (ast.Expr, *loxError.LoxError) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(token.LEFT_PAREN) {
			if val, err := p.finishCall(expr); err == nil {
				expr = val
			} else {
				return nil, err
			}
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

	return expr, nil
}

func (p *Parser) primary() (ast.Expr, *loxError.LoxError) {
	if p.match(token.FALSE) {
		return &ast.Literal{
			Value: false,
		}, nil
	}
	if p.match(token.TRUE) {
		return &ast.Literal{
			Value: true,
		}, nil
	}
	if p.match(token.NIL) {
		return &ast.Literal{
			Value: nil,
		}, nil
	}

	if p.match(token.NUMBER, token.STRING) {
		return &ast.Literal{
			Value: p.previous().Literal,
		}, nil
	}

	if p.match(token.SUPER) {
		keyword := p.previous()
		p.consume(token.DOT, "Expect '.' after 'super'.")
		method := p.consume(token.IDENTIFIER, "Expect superclass method name.")
		return &ast.Super{
			Keyword: keyword,
			Method:  method,
		}, nil
	}

	if p.match(token.THIS) {
		return &ast.This{
			Keyword: p.previous(),
		}, nil
	}

	if p.match(token.IDENTIFIER) {
		return &ast.Variable{
			Name: p.previous(),
		}, nil
	}

	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		return &ast.Grouping{
			Expression: expr,
		}, nil
	}

	// Error handling if no valid expression is found
	return nil, loxError.NewParseError(p.peek(), "Expect expression.")
}

func (p *Parser) Parse() (statements []ast.Stmt, err *loxError.LoxError) {
	// var statements []ast.Stmt
	// var err *loxError.LoxError

	defer func() {
		if r := recover(); r != nil {
			// Check if it's a parse error
			if loxErr, ok := r.(*loxError.LoxError); ok {
				err = loxErr
				return
			}
			// Re-panic for unexpected errors
			panic(r)
		}
	}()

	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
			// p.synchronize()
			// continue
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}
	/*
		if err != nil {
			return nil, err
		}
	*/
	return // statements, nil
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

func (p *Parser) consume(expected token.TokenType, message string) token.Token {
	if p.check(expected) {
		return p.advance()
	}

	err := loxError.NewParseError(p.peek(), message)
	// if err != nil {
	//	panic(err)
	// }
	loxError.ReportError(err)
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
		case token.RIGHT_BRACE: // Recover at class/method boundaries
			p.advance()
			return
		}

		p.advance()
	}
}
