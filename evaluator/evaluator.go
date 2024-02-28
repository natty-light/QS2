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
		s.Set(node.Name.Value, val, node.Constant)

	// Literals
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value, TokenLine: node.Token.Line}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value, node.Token.Line)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Scope: s, Body: body}
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

		return applyFunction(function, args)
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
	val, ok, _ := s.Get(node.Value)

	if !ok {
		return newError(node.Token.Line, "identifier not found: %s", node.Value)
	}

	return val.Value
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

// Function calls
func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)

	if !ok {
		return newError(fn.Line(), "cannot call non-function value %s", fn.Type())
	}

	exentedScope := extendFunctionScope(function, args)
	evaluated := Eval(function.Body, exentedScope)
	return unwrapReturnValue(evaluated)

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
		scope.Set(param.Value, args[paramIdx], true) // arguments from a function should be constant
	}

	return scope
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnVal, ok := obj.(*object.ReturnValue); ok {
		return returnVal.Value
	}

	return obj
}
