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
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

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

func (p *Parser) parseExpression(precedence Precedence) ast.Expr {
	prefix := p.prefixParseFns[p.currToken.Type]

	if prefix == nil {
		return nil
	}

	left := prefix()

	return left
}

// this is an prefixParseFn, so it will not call p.nextToken()
func (p *Parser) parseIdentifier() ast.Expr {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

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
