package evaluator

import (
	"fmt"
	"quonk/ast"
	"quonk/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, s *object.Scope) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, s)
	case *ast.ExpressionStmt:
		return Eval(node.Expr, s)
	case *ast.ReturnStmt:
		val := Eval(node.ReturnValue, s)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.BlockStmt:
		return evalBlockStmt(node, s)
	case *ast.VarDeclarationStmt:
		val := Eval(node.Value, s)
		if isError(val) {
			return val
		}
		s.DeclareVar(node.Name.Value, val, node.Constant, node.Token.Line)
	case *ast.VarAssignmentStmt:
		val := Eval(node.Value, s)
		if isError(val) {
			return val
		}
		errorMaybe := s.AssignVar(node.Identifier.Value, val, node.Token.Line)
		if isError(errorMaybe) {
			return errorMaybe
		}
	case *ast.ForStmt:
		errorMaybe := evalForStmt(node, s)
		if isError(errorMaybe) {
			return errorMaybe
		}
	// Literals
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Scope: s, Body: body}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, s)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.NullLiteral:
		return NULL
	case *ast.HashLiteral:
		return evalHashLiteral(node, s, node.Token.Line)
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	// Expressions
	case *ast.Identifier:
		return evalIdentifier(node, s)
	case *ast.PrefixExpr:
		right := Eval(node.Right, s)
		if isError(right) {
			return right
		}
		return evalPrefixExpr(node.Operator, right, node.Token.Line)
	case *ast.InfixExpr:
		left := Eval(node.Left, s)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, s)
		if isError(right) {
			return right
		}

		return evalInfixExpr(node.Operator, left, right, node.Token.Line)
	case *ast.IfExpr:
		return evalIfExpr(node, s)
	case *ast.CallExpr:
		if node.Function.TokenLiteral() == "quote" {
			return quote(node.Arguments[0], s)
		}
		function := Eval(node.Function, s)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, s)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args, node.Token.Line)
	case *ast.IndexExpr:
		left := Eval(node.Left, s)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, s)
		if isError(index) {
			return index
		}

		return evalIndexExpr(left, index, node.Token.Line)
	}

	return nil
}

// Statements
func evalProgram(program *ast.Program, s *object.Scope) object.Object {
	var result object.Object

	for _, stmt := range program.Stmts {
		result = Eval(stmt, s)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStmt(block *ast.BlockStmt, s *object.Scope) object.Object {
	var result object.Object

	for _, stmt := range block.Stmts {
		result = Eval(stmt, s)

		// we do not unwrap the return value here so it can bubble up
		if result != nil {
			rt := result.Type()
			if rt == object.ReturnValueObj || rt == object.ErrorObj {
				return result
			}
		}
	}

	return result
}

// Expressions
func evalPrefixExpr(operator string, right object.Object, line int) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpr(right)
	case "-":
		return evalMinusOperatorExpr(right, line)
	default:
		return newError(line, "unknown operation %s for type %s", operator, right.Type())
	}
}

func evalBangOperatorExpr(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		if right.Type() == object.IntegerObj && right.(*object.Integer).Value == 0 {
			return TRUE
		} else if right.Type() == object.FloatObj && right.(*object.Float).Value == 0 {
			return TRUE
		}

		return FALSE
	}
}

func evalMinusOperatorExpr(right object.Object, line int) object.Object {
	if right.Type() != object.IntegerObj && right.Type() != object.FloatObj {
		return newError(line, "unknown operation - for type %s", string(right.Type()))
	}

	if right.Type() == object.IntegerObj {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	} else {
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}
}

// The order of the switch statements matter here
func evalInfixExpr(operator string, left, right object.Object, line int) object.Object {
	switch {
	case left.Type() == object.IntegerObj && right.Type() == object.IntegerObj:
		return evalIntegerInfixExpr(operator, left, right, line)
	case left.Type() == object.StringObj && right.Type() == object.StringObj:
		return evalStringInfixExpr(operator, left, right, line)
	case left.Type() == object.FloatObj && right.Type() == object.FloatObj:
		return evalFloatInfixExpr(operator, left, right, line)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(right != left)
	case operator == "&&" && left.Type() == object.BooleanObj && right.Type() == object.BooleanObj:
		fallthrough
	case operator == "||" && left.Type() == object.BooleanObj && right.Type() == object.BooleanObj:
		return evalBooleanComparisonExpr(operator, left, right)
	case left.Type() != right.Type():
		return newError(line, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError(line, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpr(operator string, left, right object.Object, line int) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError(line, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpr(operator string, left, right object.Object, line int) object.Object {
	if operator != "+" {
		return newError(line, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	// Not sure about this line number here
	return &object.String{Value: leftVal + rightVal}
}

func evalBooleanComparisonExpr(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "&&":
		return nativeBoolToBooleanObject(leftVal && rightVal)
	case "||":
		return nativeBoolToBooleanObject(leftVal || rightVal)
	default:
		return NULL
	}
}

func evalIfExpr(expr *ast.IfExpr, s *object.Scope) object.Object {
	condition := Eval(expr.Condition, s)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(expr.Consequence, s)
	} else if expr.Alternative != nil {
		return Eval(expr.Alternative, s)
	} else {
		return NULL
	}
}

func evalIdentifier(node *ast.Identifier, s *object.Scope) object.Object {
	if val, _, ok := s.Get(node.Value); ok {
		return val.Value
	}

	if builtin, ok := builtIns[node.Value]; ok {
		return builtin
	}

	return newError(node.Token.Line, "identifier not found: %s", node.Value)

}

func evalExpressions(exprs []ast.Expr, s *object.Scope) []object.Object {
	var result []object.Object

	for _, e := range exprs {
		evaluated := Eval(e, s)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalIndexExpr(left, index object.Object, line int) object.Object {
	switch {
	case left.Type() == object.ArrayObj && index.Type() == object.IntegerObj:
		return evalArrayIndexExpr(left, index, line)
	case left.Type() == object.HashObj:
		return evalHashIndexExpr(left, index, line)
	default:
		return newError(line, "index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpr(array, index object.Object, line int) object.Object {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	arrLen := int64(len(arrayObj.Elements))
	maxIdx := arrLen - 1

	if (idx >= 0 && idx > maxIdx) || (idx < 0 && idx < -arrLen) {
		return newError(line, "array index out of bounds")
	}

	if idx >= 0 {
		return arrayObj.Elements[idx]
	} else {
		// since idx < 0 here, we check against the max len. Example: idx = -2, len = 3 will return elems[1],
		// the second to last elem
		return arrayObj.Elements[arrLen+idx]
	}
}

func evalForStmt(node *ast.ForStmt, s *object.Scope) object.Object {
	conditionVal := Eval(node.Condition, s)

	if conditionVal.Type() != object.BooleanObj {
		return newError(node.Token.Line, "condition for for loop must evaluate to a boolean")
	}
	condition := conditionVal.(*object.Boolean).Value

	for condition {
		Eval(node.Body, s)

		conditionVal = Eval(node.Condition, s)

		if conditionVal.Type() != object.BooleanObj {
			return newError(node.Token.Line, "condition for for loop must evaluate to a boolean")
		}
		condition = conditionVal.(*object.Boolean).Value

	}
	return nil
}

func evalHashLiteral(node *ast.HashLiteral, s *object.Scope, line int) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, s)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError(line, "unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, s)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

// Function calls
func applyFunction(fn object.Object, args []object.Object, line int) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedScope := extendFunctionScope(fn, args, line)
		evaluated := Eval(fn.Body, extendedScope)
		return unwrapReturnValue(evaluated)
	case *object.BuiltIn:
		return fn.Fn(line, args...)
	default:
		return newError(line, "not a function: %s", fn.Type())
	}
}

func evalHashIndexExpr(hash, index object.Object, line int) object.Object {
	hashObj := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError(line, "unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

// Utilty functions
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}

func newError(line int, format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...), OriginLine: line}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorObj
	}
	return false
}

func extendFunctionScope(fn *object.Function, args []object.Object, line int) *object.Scope {
	scope := object.NewEnclosedScope(fn.Scope)

	for paramIdx, param := range fn.Parameters {
		scope.DeclareVar(param.Value, args[paramIdx], true, line) // arguments from a function should be constant
	}

	return scope
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnVal, ok := obj.(*object.ReturnValue); ok {
		return returnVal.Value
	}

	return obj
}

func evalFloatInfixExpr(operator string, left, right object.Object, line int) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: int64(leftVal) % int64(rightVal)}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError(line, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}
