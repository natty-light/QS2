package compiler

import (
	"fmt"
	"quonk/ast"
	"quonk/code"
	"quonk/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) (object.ObjectType, error) {
	var err error
	var t object.ObjectType
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Stmts {
			t, err = c.Compile(s)
			if err != nil {
				return object.ErrorObj, err
			}
		}
	case *ast.ExpressionStmt:
		t, err = c.Compile(node.Expr)
		if err != nil {
			return object.ErrorObj, err
		}
		c.emit(code.OpPop)
	case *ast.InfixExpr:
		if node.Operator == "<" || node.Operator == "<=" {
			rightType, err := c.Compile(node.Right)
			if err != nil {
				return object.ErrorObj, err
			}

			leftType, err := c.Compile(node.Left)
			if err != nil {
				return object.ErrorObj, err
			}

			if leftType != rightType {
				return object.ErrorObj, fmt.Errorf("type mismatch: %s %s %s", leftType, node.Operator, rightType)
			}

			if node.Operator == "<" {
				c.emit(code.OpGt)
			} else {
				c.emit(code.OpGte)
			}
			return object.BooleanObj, nil
		}

		leftType, err := c.Compile(node.Left)
		if err != nil {
			return object.NullObj, err
		}

		rightType, err := c.Compile(node.Right)
		if err != nil {
			return object.ErrorObj, err
		}

		if leftType != rightType {
			return object.ErrorObj, fmt.Errorf("type mismatch: %s %s %s", leftType, node.Operator, rightType)
		}

		t = object.BooleanObj

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		case ">":
			c.emit(code.OpGt)
		case ">=":
			c.emit(code.OpGte)
		case "&&":
			c.emit(code.OpAnd)
		case "||":
			c.emit(code.OpOr)
		default:
			return object.ErrorObj, fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IntegerLiteral:
		t = object.IntegerObj
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	case *ast.FloatLiteral:
		t = object.FloatObj
		float := &object.Float{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(float))
	case *ast.BooleanLiteral:
		t = object.BooleanObj
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	}
	return t, nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}
