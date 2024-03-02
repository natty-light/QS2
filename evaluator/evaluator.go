package evaluator

import (
	"QuonkScript/ast"
	"QuonkScript/object"
	"fmt"
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
		s.DeclareVar(node.Name.Value, val, node.Constant)
	case *ast.VarAssignmentStmt:
		val := Eval(node.Value, s)
		if isError(val) {
			return val
		}
		errorMaybe := s.AssignVar(node.Identifier.Value, val)
		if isError(errorMaybe) {
			return errorMaybe
		}
	// Literals
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value, TokenLine: node.Token.Line}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value, node.Token.Line)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Scope: s, Body: body}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value, TokenLine: node.Token.Line}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, s)
		if len(elements) == 0 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements, TokenLine: node.Token.Line}
	case *ast.NullLiteral:
		return NULL
	// Expressions
	case *ast.Identifier:
		return evalIdentifier(node, s)
	case *ast.PrefixExpr:
		right := Eval(node.Right, s)
		if isError(right) {
			return right
		}
		return evalPrefixExpr(node.Operator, right)
	case *ast.InfixExpr:
		left := Eval(node.Left, s)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, s)
		if isError(right) {
			return right
		}

		return evalInfixExpr(node.Operator, left, right)
	case *ast.IfExpr:
		return evalIfExpr(node, s)
	case *ast.CallExpr:
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

		return evalIndexExpr(left, index)
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
func evalPrefixExpr(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpr(right)
	case "-":
		return evalMinusOperatorExpr(right)
	default:
		return newError(right.Line(), "unknown operation %s for type %s", operator, right.Type())
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
		return FALSE
	}
}

func evalMinusOperatorExpr(right object.Object) object.Object {
	if right.Type() != object.IntegerObj {
		return newError(right.Line(), "unknown operation - for type %s", string(right.Type()))
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// The order of the switch statements matter here
func evalInfixExpr(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerObj && right.Type() == object.IntegerObj:
		return evalIntegerInfixExpr(operator, left, right)
	case left.Type() == object.StringObj && right.Type() == object.StringObj:
		return evalStringInfixExpr(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right, left.Line())
	case operator == "!=":
		return nativeBoolToBooleanObject(right != left, left.Line())
	case operator == "&&" && left.Type() == object.BooleanObj && right.Type() == object.BooleanObj:
		fallthrough
	case operator == "||" && left.Type() == object.BooleanObj && right.Type() == object.BooleanObj:
		return evalBooleanComparisonExpr(operator, left, right)
	case left.Type() != right.Type():
		return newError(left.Line(), "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError(left.Line(), "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpr(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal, TokenLine: left.Line()}
	case "-":
		return &object.Integer{Value: leftVal - rightVal, TokenLine: left.Line()}
	case "*":
		return &object.Integer{Value: leftVal * rightVal, TokenLine: left.Line()}
	case "/":
		return &object.Integer{Value: leftVal / rightVal, TokenLine: left.Line()}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal, left.Line())
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal, left.Line())
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal, left.Line())
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal, left.Line())
	default:
		return newError(left.Line(), "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpr(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError(left.Line(), "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	// Not sure about this line number here
	return &object.String{Value: leftVal + rightVal, TokenLine: left.Line()}
}

func evalBooleanComparisonExpr(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "&&":
		return nativeBoolToBooleanObject(leftVal && rightVal, left.Line())
	case "||":
		return nativeBoolToBooleanObject(leftVal || rightVal, left.Line())
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

func evalIndexExpr(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ArrayObj && index.Type() == object.IntegerObj:
		return evalArrayIndexExpr(left, index)
	default:
		return newError(left.Line(), "index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpr(array, index object.Object) object.Object {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	arrLen := int64(len(arrayObj.Elements))
	maxIdx := arrLen - 1

	if (idx >= 0 && idx > maxIdx) || (idx < 0 && idx < -arrLen) {
		return newError(arrayObj.TokenLine, "array index out of bounds")
	}

	if idx >= 0 {
		return arrayObj.Elements[idx]
	} else {
		// since idx < 0 here, we check against the max len. Example: idx = -2, len = 3 will return elems[1],
		// the second to last elem
		return arrayObj.Elements[arrLen+idx]
	}
}

// Function calls
func applyFunction(fn object.Object, args []object.Object, line int) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedScope := extendFunctionScope(fn, args)
		evaluated := Eval(fn.Body, extendedScope)
		return unwrapReturnValue(evaluated)
	case *object.BuiltIn:
		return fn.Fn(line, args...)
	default:
		return newError(line, "not a function: %s", fn.Type())
	}
}

// Utilty functions
func nativeBoolToBooleanObject(input bool, line int) *object.Boolean {
	if input {
		TRUE.TokenLine = line
		return TRUE
	}
	FALSE.TokenLine = line
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

func extendFunctionScope(fn *object.Function, args []object.Object) *object.Scope {
	scope := object.NewEnclosedScope(fn.Scope)

	for paramIdx, param := range fn.Parameters {
		scope.DeclareVar(param.Value, args[paramIdx], true) // arguments from a function should be constant
	}

	return scope
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnVal, ok := obj.(*object.ReturnValue); ok {
		return returnVal.Value
	}

	return obj
}
