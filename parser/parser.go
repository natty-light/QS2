package parser

import (
	"QuonkScript/ast"
	"QuonkScript/lexer"
	"QuonkScript/token"
	"fmt"
	"strconv"
)

type (
	prefixParseFn func() ast.Expr
	infixParseFn  func(ast.Expr) ast.Expr
)

type Precedence int

const (
	LOWEST Precedence = iota + 1
	ANDOR             // I think this is right
	EQUALS
	LESSGREATEREQUAL
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]Precedence{
	token.And:                ANDOR,
	token.Or:                 ANDOR,
	token.EqualTo:            EQUALS,
	token.NotEqualTo:         EQUALS,
	token.LessThan:           LESSGREATER,
	token.GreaterThan:        LESSGREATER,
	token.GreaterThanEqualTo: LESSGREATEREQUAL,
	token.LessThanEqualTo:    LESSGREATEREQUAL,
	token.Plus:               SUM,
	token.Minus:              SUM,
	token.Slash:              PRODUCT,
	token.Star:               PRODUCT,
}

type Parser struct {
	lexer *lexer.Lexer

	currToken token.Token
	peekToken token.Token

	errors []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: make([]string, 0)}

	// peekToken and currToken are initialized to the zero value of token.Token, so we advance twice
	p.nextToken() // set peek
	p.nextToken() // set curr and peek

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)

	p.registerPrefix(token.Identifier, p.parseIdentifier)
	p.registerPrefix(token.Integer, p.parseIntegerLiteral)
	p.registerPrefix(token.Bang, p.parsePrefixExpr)
	p.registerPrefix(token.Minus, p.parsePrefixExpr)
	p.registerPrefix(token.True, p.parseBooleanLiteral)
	p.registerPrefix(token.False, p.parseBooleanLiteral)
	p.registerPrefix(token.LeftParen, p.parseGroupedExpr)
	p.registerPrefix(token.If, p.parseIfExpr)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.Plus, p.parseInfixExpr)
	p.registerInfix(token.Minus, p.parseInfixExpr)
	p.registerInfix(token.Slash, p.parseInfixExpr)
	p.registerInfix(token.Star, p.parseInfixExpr)
	p.registerInfix(token.EqualTo, p.parseInfixExpr)
	p.registerInfix(token.NotEqualTo, p.parseInfixExpr)
	p.registerInfix(token.GreaterThanEqualTo, p.parseInfixExpr)
	p.registerInfix(token.LessThanEqualTo, p.parseInfixExpr)
	p.registerInfix(token.GreaterThan, p.parseInfixExpr)
	p.registerInfix(token.LessThan, p.parseInfixExpr)
	p.registerInfix(token.And, p.parseInfixExpr)
	p.registerInfix(token.Or, p.parseInfixExpr)

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// advances current and peek by one
func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// Checks whether current token matches given type
func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

// checks whether peek token matches given type
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// Checks if peek token matches given type, advances tokens if true
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken() // eats
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead on line %d",
		t, p.peekToken.Type, p.peekToken.Line)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() Precedence {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currPrecedence() Precedence {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Parsing methods
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Stmts = make([]ast.Stmt, 0)

	for !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Stmts = append(program.Stmts, stmt)
		}
		p.nextToken() // advance past semis?
	}
	return program
}

// Statements
func (p *Parser) parseStatement() ast.Stmt {
	switch p.currToken.Type {
	case token.Mut:
		fallthrough
	case token.Const:
		return p.parseVarDeclarationStmt()
	case token.Return:
		return p.parseReturnStmt()
	default:
		return p.parseExpressionStmt()
	}
}

func (p *Parser) parseVarDeclarationStmt() *ast.VarDeclarationStmt {

	// To be here, currToken is either Mut or Const
	isConst := p.currToken.Type == token.Const

	stmt := &ast.VarDeclarationStmt{Token: p.currToken, Constant: isConst}

	// expectPeek eats?
	if !p.expectPeek(token.Identifier) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.Assign) {
		return nil
	}

	for !p.currTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{Token: p.currToken}

	p.nextToken()

	for !p.currTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStmt() *ast.ExpressionStmt {
	stmt := &ast.ExpressionStmt{Token: p.currToken}

	stmt.Expr = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	block := &ast.BlockStmt{Token: p.currToken}
	block.Stmts = make([]ast.Stmt, 0)

	p.nextToken() // advance past {

	for !p.currTokenIs(token.RightCurlyBracket) && !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Stmts = append(block.Stmts, stmt)
		}
		p.nextToken() // what is this advance for
	}
	return block
}

// Expressions
func (p *Parser) parseExpression(precedence Precedence) ast.Expr {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}

	left := prefix()

	// if the statement has not ended and the passed in precedence is lower than the precedence of the next token
	for !p.peekTokenIs(token.Semicolon) && precedence < p.peekPrecedence() {
		// look for an infix parse fn
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return left
		}

		p.nextToken()

		// this bubbles up
		left = infix(left)
	}

	return left
}

// prefix and infix functions

// this is an prefixParseFn, so it will not call p.nextToken() at the end
func (p *Parser) parseIdentifier() ast.Expr {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

// this is an prefixParseFn, so it will not call p.nextToken() at the end
func (p *Parser) parseIntegerLiteral() ast.Expr {
	literal := &ast.IntegerLiteral{Token: p.currToken}

	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	literal.Value = value

	return literal
}

// this is an prefixParseFn, so it will not call p.nextToken() at the end
func (p *Parser) parsePrefixExpr() ast.Expr {
	expr := &ast.PrefixExpr{Token: p.currToken, Operator: p.currToken.Literal}

	p.nextToken() // advance past operator
	expr.Right = p.parseExpression(PREFIX)

	return expr
}

// this is an infixParseFn, so it will not call p.nextToken() at the end
func (p *Parser) parseInfixExpr(left ast.Expr) ast.Expr {
	expr := &ast.InfixExpr{Token: p.currToken, Operator: p.currToken.Literal, Left: left}

	precedence := p.currPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseBooleanLiteral() ast.Expr {
	return &ast.BooleanLiteral{Token: p.currToken, Value: p.currTokenIs(token.True)}
}

func (p *Parser) parseGroupedExpr() ast.Expr {
	p.nextToken() // advance past (

	expr := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RightParen) {
		return nil
	}

	return expr
}

func (p *Parser) parseIfExpr() ast.Expr {
	expr := &ast.IfExpr{Token: p.currToken}

	if !p.expectPeek(token.LeftParen) {
		return nil
	}
	p.nextToken() // advance past (
	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RightParen) {
		return nil
	}
	if !p.expectPeek(token.LeftCurlyBracket) {
		return nil
	}

	expr.Consequence = p.parseBlockStmt()

	// For else if, probably want to set the alternative to another ast.IfExpr, maybe by calling this method recursively
	if p.peekTokenIs(token.Else) {
		p.nextToken() // advance past else
		if !p.expectPeek(token.LeftCurlyBracket) {
			return nil
		}

		expr.Alternative = p.parseBlockStmt()
	}

	return expr
}
