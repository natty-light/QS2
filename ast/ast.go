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

type (
	Program struct {
		Stmts []Stmt
	}

	VarDeclarationStmt struct {
		Token    token.Token // token.Mut or token.Const
		Name     *Identifier
		Value    Expr
		Constant bool
	}

	Identifier struct {
		Token token.Token // token.Ident
		Value string
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

// Node interfaces
func (p *Program) TokenLiteral() string {
	if len(p.Stmts) > 0 {
		return p.Stmts[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, stmt := range p.Stmts {
		out.WriteString(stmt.String())
	}

	return out.String()
}

func (v *VarDeclarationStmt) TokenLiteral() string {
	return v.Token.Literal
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

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

func (r *ReturnStmt) TokenLiteral() string {
	return r.Token.Literal
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

func (e *ExpressionStmt) TokenLiteral() string {
	return e.Token.Literal
}

func (e *ExpressionStmt) String() string {
	if e.Expr != nil {
		return e.Expr.String()
	}
	return ""
}

// Statements
func (v *VarDeclarationStmt) statementNode() {}
func (r *ReturnStmt) statementNode()         {}
func (e *ExpressionStmt) statementNode()     {}

// Expressions
func (i *Identifier) expressionNode() {}
