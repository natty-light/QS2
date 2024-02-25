package ast

import (
	"QuonkScript/token"
	"bytes"
)

// Interfaces

type (
	Node interface {
		TokenLiteral() string
		String() string
	}

	Stmt interface {
		Node
		statementNode()
	}

	Expr interface {
		Node
		expressionNode()
	}
)

// Node
type (
	Program struct {
		Stmts []Stmt
	}
)

// Statements
type (
	VarDeclarationStmt struct {
		Token    token.Token // token.Mut or token.Const
		Name     *Identifier
		Value    Expr
		Constant bool
	}

	ReturnStmt struct {
		Token       token.Token
		ReturnValue Expr
	}

	ExpressionStmt struct {
		Token token.Token
		Expr  Expr
	}
)

// Expressions and literals
type (
	Identifier struct {
		Token token.Token // token.Ident
		Value string
	}

	IntegerLiteral struct {
		Token token.Token
		Value int64
	}

	PrefixExpr struct {
		Token    token.Token
		Operator string
		Right    Expr
	}

	InfixExpr struct {
		Token    token.Token
		Left     Expr
		Operator string
		Right    Expr
	}
)

// Node interfaces
func (p *Program) TokenLiteral() string {
	if len(p.Stmts) > 0 {
		return p.Stmts[0].TokenLiteral()
	} else {
		return ""
	}
}

func (v *VarDeclarationStmt) TokenLiteral() string {
	return v.Token.Literal
}

func (r *ReturnStmt) TokenLiteral() string {
	return r.Token.Literal
}

func (e *ExpressionStmt) TokenLiteral() string {
	return e.Token.Literal
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (p *PrefixExpr) TokenLiteral() string {
	return p.Token.Literal
}

func (i *InfixExpr) TokenLiteral() string {
	return i.Token.Literal
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, stmt := range p.Stmts {
		out.WriteString(stmt.String())
	}

	return out.String()
}

func (v *VarDeclarationStmt) String() string {
	var out bytes.Buffer

	out.WriteString(v.TokenLiteral() + " ")
	out.WriteString(v.Name.String())
	out.WriteString(" = ")

	if v.Value != nil {
		out.WriteString(v.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

func (r *ReturnStmt) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral() + " ")

	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

func (e *ExpressionStmt) String() string {
	if e.Expr != nil {
		return e.Expr.String()
	}
	return ""
}

// Expressions
func (i *Identifier) String() string {
	return i.Value
}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

func (p *PrefixExpr) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

func (i *InfixExpr) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}

// Statements
func (v *VarDeclarationStmt) statementNode() {}
func (r *ReturnStmt) statementNode()         {}
func (e *ExpressionStmt) statementNode()     {}

// Expressions
func (i *Identifier) expressionNode()     {}
func (i *IntegerLiteral) expressionNode() {}
func (p *PrefixExpr) expressionNode()     {}
func (i *InfixExpr) expressionNode()      {}
