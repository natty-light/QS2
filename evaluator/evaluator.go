package evaluator

import (
	"QuonkScript/ast"
	"QuonkScript/object"
)

// I am not sure why I have to do these casts and the book doesn't
func Eval(node ast.Node) object.Object {
	switch node.(type) {
	case *ast.Program:
		return evalStatements(node.(*ast.Program).Stmts)
	case *ast.ExpressionStmt:
		return Eval(node.(*ast.ExpressionStmt).Expr)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.(*ast.IntegerLiteral).Value}
	}

	return nil
}

func evalStatements(stmts []ast.Stmt) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}
