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

	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

func New() *Compiler {
	return &Compiler{
		instructions:        code.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
}

func (c *Compiler) Compile(node ast.Node) (object.ObjectType, error) {
	var err error
	var t object.ObjectType
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Stmts {
			_, err = c.Compile(s)
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
	case *ast.BlockStmt:
		for _, s := range node.Stmts {
			_, err = c.Compile(s)
			if err != nil {
				return object.ErrorObj, err
			}
		}
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
				return object.ErrorObj, fmt.Errorf("type mismatch: %s %s %s on line %d", leftType, node.Operator, rightType, node.Token.Line)
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
			return object.ErrorObj, fmt.Errorf("type mismatch: %s %s %s on line %d", leftType, node.Operator, rightType, node.Token.Line)
		}

		switch node.Operator {

		case "+":
			t = leftType
			c.emit(code.OpAdd)
		case "-":
			t = leftType
			c.emit(code.OpSub)
		case "*":
			t = leftType
			c.emit(code.OpMul)
		case "/":
			t = leftType
			c.emit(code.OpDiv)
		case "==":
			t = object.BooleanObj
			c.emit(code.OpEqual)
		case "!=":
			t = object.BooleanObj
			c.emit(code.OpNotEqual)
		case ">":
			t = object.BooleanObj
			c.emit(code.OpGt)
		case ">=":
			t = object.BooleanObj
			c.emit(code.OpGte)
		case "&&":
			t = object.BooleanObj
			c.emit(code.OpAnd)
			t = object.BooleanObj
		case "||":
			c.emit(code.OpOr)
		default:
			return object.ErrorObj, fmt.Errorf("unknown operator %s on line %d", node.Operator, node.Token.Line)
		}
	case *ast.PrefixExpr:
		t, err = c.Compile(node.Right)
		if err != nil {
			return object.ErrorObj, err
		}

		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			if t != object.IntegerObj && t != object.FloatObj {
				return object.ErrorObj, fmt.Errorf("unknown operator %s for type %s on line %d", node.Operator, t, node.Token.Line)
			}
			c.emit(code.OpMinus)
		default:
			return object.ErrorObj, fmt.Errorf("unknown operator %s on line %d", node.Operator, node.Token.Line)
		}
	case *ast.IfExpr:
		// we don't need to update t here because we're not bubbling the value back up like in expressions
		_, err = c.Compile(node.Condition)
		if err != nil {
			return object.ErrorObj, err
		}
		// emit with operand to be replaced later
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		_, err = c.Compile(node.Consequence)
		if err != nil {
			return object.ErrorObj, err
		}

		// remove last pop after compiling consequence so we don't inadvertently pop too many times
		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}

		// only if there is no alternative do we jump to immediately after the consequence
		if node.Alternative == nil {
			afterConsequencePos := len(c.instructions)
			c.changeOperand(jumpNotTruthyPos, afterConsequencePos)
		} else {
			// emit with operand to be replaced later
			jumpPos := c.emit(code.OpJump, 9999)

			afterConsequencePos := len(c.instructions)
			c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

			_, err = c.Compile(node.Alternative)
			if err != nil {
				return object.ErrorObj, err
			}

			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}

			afterAlternativePos := len(c.instructions)
			c.changeOperand(jumpPos, afterAlternativePos)
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

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	return pos
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.previousInstruction = previous
	c.lastInstruction = last
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == code.OpPop
}

func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.instructions[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.instructions[opPos])
	newInstruction := code.Make(op, operand)
	c.replaceInstruction(opPos, newInstruction)
}
