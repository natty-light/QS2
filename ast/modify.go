package ast

type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {
	case *Program:
		for i, stmt := range node.Stmts {
			node.Stmts[i], _ = Modify(stmt, modifier).(Stmt)
		}
	case *ExpressionStmt:
		node.Expr, _ = Modify(node.Expr, modifier).(Expr)
	case *InfixExpr:
		node.Left, _ = Modify(node.Left, modifier).(Expr)
		node.Right, _ = Modify(node.Right, modifier).(Expr)
	case *PrefixExpr:
		node.Right, _ = Modify(node.Right, modifier).(Expr)
	case *IndexExpr:
		node.Left, _ = Modify(node.Left, modifier).(Expr)
		node.Index, _ = Modify(node.Index, modifier).(Expr)
	case *IfExpr:
		node.Condition, _ = Modify(node.Condition, modifier).(Expr)
		node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStmt)
		if node.Alternative != nil {
			node.Alternative, _ = Modify(node.Alternative, modifier).(*BlockStmt)
		}
	case *BlockStmt:
		for i, _ := range node.Stmts {
			node.Stmts[i], _ = Modify(node.Stmts[i], modifier).(Stmt)
		}
	case *ReturnStmt:
		node.ReturnValue, _ = Modify(node.ReturnValue, modifier).(Expr)
	case *VarDeclarationStmt:
		node.Value, _ = Modify(node.Value, modifier).(Expr)
	case *VarAssignmentStmt:
		node.Value, _ = Modify(node.Value, modifier).(Expr)
	case *ForStmt:
		node.Condition, _ = Modify(node.Condition, modifier).(Expr)
		node.Body, _ = Modify(node.Body, modifier).(*BlockStmt)
	case *FunctionLiteral:
		for i, _ := range node.Parameters {
			node.Parameters[i], _ = Modify(node.Parameters[i], modifier).(*Identifier)
		}
		node.Body, _ = Modify(node.Body, modifier).(*BlockStmt)
	case *ArrayLiteral:
		for i, _ := range node.Elements {
			node.Elements[i], _ = Modify(node.Elements[i], modifier).(Expr)
		}
	case *HashLiteral:
		newPairs := make(map[Expr]Expr)
		for key, val := range node.Pairs {
			newKey, _ := Modify(key, modifier).(Expr)
			newVal, _ := Modify(val, modifier).(Expr)
			newPairs[newKey] = newVal
		}
		node.Pairs = newPairs
	}
	return modifier(node)
}
