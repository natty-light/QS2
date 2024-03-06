package ast

import (
	"reflect"
	"testing"
)

func TestModify(t *testing.T) {
	one := func() Expr { return &IntegerLiteral{Value: 1} }
	two := func() Expr { return &IntegerLiteral{Value: 2} }

	turnOneIntoTwo := func(node Node) Node {
		integer, ok := node.(*IntegerLiteral)
		if !ok {
			return node
		}

		if integer.Value != 1 {
			return node
		}

		integer.Value = 2
		return integer
	}

	tests := []struct {
		input    Node
		expected Node
	}{
		{
			one(),
			two(),
		},
		{
			&Program{
				Stmts: []Stmt{
					&ExpressionStmt{
						Expr: one(),
					},
				},
			},
			&Program{
				Stmts: []Stmt{
					&ExpressionStmt{
						Expr: two(),
					},
				},
			},
		},
		{
			&InfixExpr{Left: one(), Operator: "+", Right: two()},
			&InfixExpr{Left: two(), Operator: "+", Right: two()},
		},
		{
			&InfixExpr{Left: two(), Operator: "+", Right: one()},
			&InfixExpr{Left: two(), Operator: "+", Right: two()},
		},
		{
			&PrefixExpr{Operator: "-", Right: one()},
			&PrefixExpr{Operator: "-", Right: two()},
		},
		{
			&IndexExpr{Left: one(), Index: one()},
			&IndexExpr{Left: two(), Index: two()},
		},
		{
			&IfExpr{
				Condition: one(),
				Consequence: &BlockStmt{
					Stmts: []Stmt{
						&ExpressionStmt{Expr: one()},
					},
				},
				Alternative: &BlockStmt{
					Stmts: []Stmt{
						&ExpressionStmt{Expr: one()},
					},
				},
			},
			&IfExpr{
				Condition: two(),
				Consequence: &BlockStmt{
					Stmts: []Stmt{
						&ExpressionStmt{Expr: two()},
					},
				},
				Alternative: &BlockStmt{
					Stmts: []Stmt{
						&ExpressionStmt{
							Expr: two(),
						},
					},
				},
			},
		},
		{
			&ReturnStmt{ReturnValue: one()},
			&ReturnStmt{ReturnValue: two()},
		},
		{
			&VarDeclarationStmt{Value: one()},
			&VarDeclarationStmt{Value: two()},
		},
		{
			&VarAssignmentStmt{Value: one()},
			&VarAssignmentStmt{Value: two()},
		},
		{
			&ForStmt{
				Condition: one(),
				Body: &BlockStmt{
					Stmts: []Stmt{
						&ExpressionStmt{Expr: one()},
					},
				},
			},
			&ForStmt{
				Condition: two(),
				Body: &BlockStmt{
					Stmts: []Stmt{
						&ExpressionStmt{Expr: two()},
					},
				},
			},
		},
		{
			&FunctionLiteral{
				Parameters: []*Identifier{},
				Body: &BlockStmt{
					Stmts: []Stmt{
						&ExpressionStmt{Expr: one()},
					},
				},
			},
			&FunctionLiteral{
				Parameters: []*Identifier{},
				Body: &BlockStmt{
					Stmts: []Stmt{
						&ExpressionStmt{Expr: two()},
					},
				},
			},
		},
		{
			&ArrayLiteral{
				Elements: []Expr{
					one(),
				},
			},
			&ArrayLiteral{
				Elements: []Expr{
					two(),
				},
			},
		},
	}

	for _, tt := range tests {
		modified := Modify(tt.input, turnOneIntoTwo)

		equal := reflect.DeepEqual(modified, tt.expected)
		if !equal {
			t.Errorf("not equal. got=%#v, want=%#v", modified, tt.expected)
		}
	}

	hashLiteral := &HashLiteral{
		Pairs: map[Expr]Expr{
			one(): one(),
			one(): one(),
		},
	}

	Modify(hashLiteral, turnOneIntoTwo)

	for key, val := range hashLiteral.Pairs {
		key, _ := key.(*IntegerLiteral)
		if key.Value != 2 {
			t.Errorf("value is not %d, got=%d", 2, key.Value)
		}

		val, _ := val.(*IntegerLiteral)
		if val.Value != 2 {
			t.Errorf("value is not %d, got=%d", 2, val.Value)
		}
	}
}
