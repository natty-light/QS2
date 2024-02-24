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
