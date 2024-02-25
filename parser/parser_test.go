package parser

import (
	"QuonkScript/ast"
	"QuonkScript/lexer"
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

	literal, ok := stmt.Expr.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("expr not *ast.IntegerLiteral. got=%T", stmt.Expr)
	}

	if literal.Value != 5 {
		t.Errorf("ident.Value not %d. got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}
