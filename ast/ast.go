package ast

import "QuonkScript/token"

// Interfaces

type (
	Node interface {
		TokenLiteral() string
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

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// Statements
func (v *VarDeclarationStmt) statementNode() {}

// Expressions
func (i *Identifier) expressionNode() {}
