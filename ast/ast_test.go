package ast

import (
	"quonk/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Stmts: []Stmt{
			&VarDeclarationStmt{
				Token:    token.Token{Type: token.Mut, Literal: "mut"},
				Name:     &Identifier{Token: token.Token{Type: token.Identifier, Literal: "x"}, Value: "x"},
				VarType:  &TypeLiteral{Token: token.Token{Type: token.IntType, Literal: "int"}},
				Value:    &Identifier{Token: token.Token{Type: token.Identifier, Literal: "y"}, Value: "y"},
				Constant: false,
			},
		},
	}

	if program.String() != "mut x int = y;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
