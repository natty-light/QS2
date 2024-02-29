package ast

import (
	"QuonkScript/token"
	"bytes"
	"strings"
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

	BlockStmt struct {
		Token token.Token
		Stmts []Stmt
	}

	VarAssignmentStmt struct {
		Token      token.Token
		Identifier *Identifier
		Value      Expr
	}
)

// Expressions and literals
type (
	// Literals
	IntegerLiteral struct {
		Token token.Token
		Value int64
	}

	BooleanLiteral struct {
		Token token.Token
		Value bool
	}

	FunctionLiteral struct {
		Token      token.Token
		Parameters []*Identifier
		Body       *BlockStmt
	}
	// Expressions
	Identifier struct {
		Token token.Token // token.Ident
		Value string
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

	IfExpr struct {
		Token       token.Token
		Condition   Expr
		Consequence *BlockStmt
		Alternative *BlockStmt
	}

	CallExpr struct {
		Token     token.Token
		Function  Expr
		Arguments []Expr
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

func (b *BooleanLiteral) TokenLiteral() string {
	return b.Token.Literal
}

func (i IfExpr) TokenLiteral() string {
	return i.Token.Literal
}

func (b *BlockStmt) TokenLiteral() string {
	return b.Token.Literal
}

func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}

func (c *CallExpr) TokenLiteral() string {
	return c.Token.Literal
}

func (v *VarAssignmentStmt) TokenLiteral() string {
	return v.Token.Literal
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

func (i *IfExpr) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())

	if i.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(i.Alternative.String())
	}

	return out.String()
}

func (b *BlockStmt) String() string {
	var out bytes.Buffer

	for _, stmt := range b.Stmts {
		out.WriteString(stmt.String())
	}

	return out.String()
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

func (b *BooleanLiteral) String() string {
	return b.Token.Literal
}

func (f *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := make([]string, 0)

	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(f.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())

	return out.String()
}

func (c *CallExpr) String() string {
	var out bytes.Buffer
	args := make([]string, 0)
	for _, arg := range c.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

func (v *VarAssignmentStmt) String() string {
	var out bytes.Buffer

	out.WriteString(v.Identifier.String())
	out.WriteString(" = ")
	out.WriteString(v.Value.String())

	return out.String()
}

// Statements
func (v *VarDeclarationStmt) statementNode() {}
func (r *ReturnStmt) statementNode()         {}
func (e *ExpressionStmt) statementNode()     {}
func (b *BlockStmt) statementNode()          {}
func (v *VarAssignmentStmt) statementNode()  {}

// Expressions
func (i *Identifier) expressionNode()      {}
func (i *IntegerLiteral) expressionNode()  {}
func (p *PrefixExpr) expressionNode()      {}
func (i *InfixExpr) expressionNode()       {}
func (b *BooleanLiteral) expressionNode()  {}
func (i *IfExpr) expressionNode()          {}
func (f *FunctionLiteral) expressionNode() {}
func (c *CallExpr) expressionNode()        {}
