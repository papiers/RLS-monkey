package vm

import (
	"fmt"

	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int // 始终指向栈中下一个空闲位置，栈顶元素为 stack[sp-1]
}

// New 创建一个新的虚拟机
func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

// StackTop 返回栈顶元素
func (vm *VM) StackTop() object.Object {
	if vm.sp > 0 {
		return vm.stack[vm.sp-1]
	}
	return nil
}

// Run 执行字节码
func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()
			leftVal := left.(*object.Integer).Value
			rightVal := right.(*object.Integer).Value
			err := vm.push(&object.Integer{Value: leftVal + rightVal})
			if err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		default:
			return fmt.Errorf("unknown opcode: %d", op)
		}
	}
	return nil
}

// push 将对象压入栈
func (vm *VM) push(obj object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = obj
	vm.sp++
	return nil
}

// pop 从栈中弹出对象
func (vm *VM) pop() object.Object {
	obj := vm.stack[vm.sp-1]
	vm.sp--
	return obj
}

// LastPoppedStackElem 返回最近弹出的栈元素
func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}
