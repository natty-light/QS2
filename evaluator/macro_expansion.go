package evaluator

import (
	"QuonkScript/ast"
	"QuonkScript/object"
)

func DefineMacros(program *ast.Program, scope *object.Scope) {
	definitions := []int{}

	for i, statement := range program.Stmts {
		if isMacroDefinition(statement) {
			addMacro(statement, scope)
			definitions = append(definitions, i)
		}
	}

	for i := len(definitions) - 1; i >= 0; i = i - 1 {
		definitionIndex := definitions[i]
		program.Stmts = append(program.Stmts[:definitionIndex], program.Stmts[definitionIndex+1:]...)
	}
}

// TODO: Update this to allow for things like
// let myMacro = macro(x) { x };
// let anotherName = myMacro;
func isMacroDefinition(node ast.Stmt) bool {
	switch node := node.(type) {
	case *ast.VarDeclarationStmt:
		_, ok := node.Value.(*ast.MacroLiteral)
		if !ok {
			return false
		}
		return true
	case *ast.VarAssignmentStmt:
		_, ok := node.Value.(*ast.MacroLiteral)
		if !ok {
			return false
		}
		return true
	default:
		return false
	}
}

func addMacro(stmt ast.Stmt, scope *object.Scope) {
	switch stmt := stmt.(type) {
	case *ast.VarDeclarationStmt:
		macro := stmt.Value.(*ast.MacroLiteral)
		macroObj := &object.Macro{Parameters: macro.Parameters, Body: macro.Body, Scope: scope}
		scope.Set(stmt.Name.Value, macroObj, stmt.Constant)
	case *ast.VarAssignmentStmt:
		macro := stmt.Value.(*ast.MacroLiteral)
		macroObj := &object.Macro{Parameters: macro.Parameters, Body: macro.Body, Scope: scope}
		scope.Set(stmt.Identifier.Value, macroObj, false)
	}
}
