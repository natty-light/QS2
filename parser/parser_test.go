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

	if !testIdentifier(t, stmt.Expr, "myVar") {
		return
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

	if !testIntegerLiteral(t, stmt.Expr, 5) {
		return
	}
}

func TestParsingPrefixExpr(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
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
		if !testLiteralExpr(t, expr.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true && false", true, "&&", false},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
		{"true || false", true, "||", false},
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

		if !testInfixExpr(t, stmt.Expr, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"true && false",
			"(true && false)",
		},
		{
			"true || false",
			"(true || false)",
		},
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
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
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

func TestBooleanExpr(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Stmts) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Stmts))
		}

		stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
				program.Stmts[0])
		}

		if !testLiteralExpr(t, stmt.Expr, tt.expectedBoolean) {
			return
		}
	}
}

// Utilities

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

func testIdentifier(t *testing.T, expr ast.Expr, value string) bool {
	ident, ok := expr.(*ast.Identifier)

	if !ok {
		t.Errorf("expr not *ast.Identifier. got=%T", expr)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpr(t *testing.T, expr ast.Expr, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case int64:
		return testIntegerLiteral(t, expr, v)
	case string:
		return testIdentifier(t, expr, v)
	case bool:
		return testBooleanLiteral(t, expr, v)
	}
	t.Errorf("type of expr not handled. got=%T", expr)
	return false
}

func testInfixExpr(t *testing.T, expr ast.Expr, left interface{}, operator string, right interface{}) bool {
	opExpr, ok := expr.(*ast.InfixExpr)
	if !ok {
		t.Errorf("expr is not *ast.InfixExpr. got=%T(%s)", expr, expr)
		return false
	}

	if !testLiteralExpr(t, opExpr.Left, left) {
		return false
	}

	if opExpr.Operator != operator {
		t.Errorf("expr.Operator is not '%s'. got=%s", operator, opExpr.Operator)
		return false
	}

	if !testLiteralExpr(t, opExpr.Right, right) {
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, expr ast.Expr, value bool) bool {
	boolean, ok := expr.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("expr not *ast.BooleanLiteral. got=%T", expr)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value not %t. got=%t", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral not %t. got=%s", value, boolean.TokenLiteral())
		return false
	}

	return true
}
