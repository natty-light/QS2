package vm

import (
	"fmt"
	"quonk/code"
	"quonk/compiler"
	"quonk/object"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int // always points to next value. top of stack is stack[sp - 1]
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			// get constant index from instruction
			constIdx := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2 // update instruction pointer

			// push constant onto stack
			err := vm.push(vm.constants[constIdx])
			if err != nil {
				return err
			}
		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()

			switch true {
			case left.Type() == object.IntegerObj && right.Type() == object.IntegerObj:
				leftVal := left.(*object.Integer).Value
				rightVal := right.(*object.Integer).Value
				result := leftVal + rightVal
				err := vm.push(&object.Integer{Value: result})
				if err != nil {
					panic(err)
				}
			case left.Type() == object.FloatObj && right.Type() == object.FloatObj:
				leftVal := left.(*object.Float).Value
				rightVal := right.(*object.Float).Value
				result := leftVal + rightVal
				err := vm.push(&object.Float{Value: result})
				if err != nil {
					panic(err)
				}
			default:
				return fmt.Errorf("type mismatch: %s + %s", left.Type(), right.Type())
			}
		case code.OpPop:
			vm.pop()
		}
	}
	return nil
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.head()
	vm.sp--
	return o
}

func (vm *VM) head() object.Object {
	return vm.stack[vm.sp-1]
}
