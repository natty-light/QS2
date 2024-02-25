package parser

import (
	"QuonkScript/ast"
	"QuonkScript/lexer"
	"fmt"
	"testing"
)

func TestVarDeclarationStmts(t *testing.T) {
	source := `
	mut x = 5;
	const y = 10;
	mut val = 838383;	
	`

	lexer := lexer.New(source)

	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if program == nil {
		t.Fatal("ParseProgram returned nil")
	}

	if len(program.Stmts) != 3 {
		t.Fatalf("program.Stmts does not contain 3 statements. got=%d", len(program.Stmts))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"val"},
	}

	for i, tt := range tests {
		stmt := program.Stmts[i]
		if !testVarDeclarationStmt(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testVarDeclarationStmt(t *testing.T, s ast.Stmt, name string) bool {
	if s.TokenLiteral() != "mut" && s.TokenLiteral() != "const" {
		t.Errorf("s.TokenLiteral not `mut` or `const`. got=%q", s.TokenLiteral())
		return false
	}

	varDeclStmt, ok := s.(*ast.VarDeclarationStmt)
	if !ok {
		t.Errorf("s not *ast.VarDeclarationStmt. got=%T", s)
		return false
	}

	if varDeclStmt.Name.Value != name {
		t.Errorf("varDeclStmt.Name.Value not '%s'. got=%s", name, varDeclStmt.Name.Value)
		return false
	}

	if varDeclStmt.Name.TokenLiteral() != name {
		t.Errorf("varDeclStmt.Name.TokenLiteral() not '%s'. got=%s", name, varDeclStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	source := `
	return 5;
	return 10;
	return 15;
	`

	lexer := lexer.New(source)

	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Stmts) != 3 {
		t.Fatalf("program.Stmts does not contain 3 statements. got=%d", len(program.Stmts))
	}

	for _, stmt := range program.Stmts {
		returnStmt, ok := stmt.(*ast.ReturnStmt)

		if !ok {
			t.Errorf("stmt not *ast.ReturnStmt. got=%T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return'. got=%q", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpr(t *testing.T) {
	source := "myVar"

	l := lexer.New(source)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("program.Stmts[0] not ast.ExpressionStmt. got=%T", stmt)
	}

	ident, ok := stmt.Expr.(*ast.Identifier)

	if !ok {
		t.Errorf("expr not *ast.Identifier. got=%T", stmt.Expr)
	}

	if ident.Value != "myVar" {
		t.Errorf("ident.Value not %s. got=%s", "myVar", ident.Value)
	}

	if ident.TokenLiteral() != "myVar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "myVar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpr(t *testing.T) {
	source := "5"

	l := lexer.New(source)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("program.Stmts[0] not ast.ExpressionStmt. got=%T", stmt)
	}

	testIntegerLiteral(t, stmt.Expr, 5)
}

func TestParsingPrefixExpr(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Stmts) != 1 {
			t.Fatalf("program.Stmts does not contain %d statements. got=%d\n", 1, len(program.Stmts))
		}

		stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt. got=%T", program.Stmts[0])
		}

		expr, ok := stmt.Expr.(*ast.PrefixExpr)
		if !ok {
			t.Fatalf("stmt is not *ast.PrefixExpr. got=%T", stmt.Expr)
		}

		if expr.Operator != tt.operator {
			t.Fatalf("expr.Operator is not '%s'. got=%s", tt.operator, expr.Operator)
		}
		if !testIntegerLiteral(t, expr.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, i ast.Expr, value int64) bool {
	integ, ok := i.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("i not *ast.IntegerLiteral. got=%T", integ)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Stmts) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Stmts))
		}

		stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
				program.Stmts[0])
		}

		exp, ok := stmt.Expr.(*ast.InfixExpr)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expr)
		}
		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecendeParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		}, {
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 >= 4 * 5 == 5 <= 3 - 2",
			"((3 >= (4 * 5)) == (5 <= (3 - 2)))",
		},
		{
			"3 >= 4 * 5 != 5 < 3 / 2 - 4",
			"((3 >= (4 * 5)) != (5 < ((3 / 2) - 4)))",
		},
		{
			"4 > 5 || 2 < 3",
			"((4 > 5) || (2 < 3))",
		},
		{
			"4 > 5 && 2 < 3",
			"((4 > 5) && (2 < 3))",
		},
		{
			"4 > 5 || 2 < 3 && 2 + 4 * 3 / 7",
			"(((4 > 5) || (2 < 3)) && (2 + ((4 * 3) / 7)))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
