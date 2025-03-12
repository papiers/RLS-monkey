package vm

import (
	"fmt"

	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

const (
	StackSize   = 2048
	GlobalsSize = 65535
)

var (
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
	Null  = &object.Null{}
)

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack   []object.Object
	sp      int // 始终指向栈中下一个空闲位置，栈顶元素为 stack[sp-1]
	globals []object.Object
}

// New 创建一个新的虚拟机
func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
		globals:      make([]object.Object, GlobalsSize),
	}
}

// NewWithGlobalsStore 创建一个新的虚拟机，并允许自定义全局变量存储
func NewWithGlobalsStore(bytecode *compiler.Bytecode, globals []object.Object) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
		globals:      globals,
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
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}
		case code.OpBang:
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}
		case code.OpMinus:
			err := vm.executeMinusOperator()
			if err != nil {
				return err
			}
		case code.OpJump:
			pos := code.ReadUint16(vm.instructions[ip+1:])
			ip = int(pos - 1)
		case code.OpJumpNotTruthy:
			pos := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			condition := vm.pop()
			if !isTruthy(condition) {
				ip = int(pos - 1)
			}
		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			vm.globals[globalIndex] = vm.pop()
		case code.OpGetGlobal:
			index := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.globals[index])
			if err != nil {
				return err
			}
		case code.OpArray:
			arrLen := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2
			elements := vm.buildArray(vm.sp-arrLen, vm.sp)
			vm.sp = vm.sp - arrLen
			err := vm.push(elements)
			if err != nil {
				return err
			}
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

// executeBinaryOperation 执行二元操作
func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	leftType := left.Type()
	rightType := right.Type()
	switch {
	case leftType == object.INTEGER && rightType == object.INTEGER:
		return vm.executeBinaryIntegerOperation(op, left, right)
	case leftType == object.STRING && rightType == object.STRING:
		return vm.executeBinaryStringOperation(op, left, right)
	}
	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

// executeBinaryIntegerOperation 执行二元整数操作
func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) error {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	var result int64
	switch op {
	case code.OpAdd:
		result = leftVal + rightVal
	case code.OpSub:
		result = leftVal - rightVal
	case code.OpMul:
		result = leftVal * rightVal
	case code.OpDiv:
		result = leftVal / rightVal
	default:
		return fmt.Errorf("unknown operator: %c", op)
	}
	return vm.push(&object.Integer{Value: result})
}

// executeBinaryStringOperation 执行二元字符串操作
func (vm *VM) executeBinaryStringOperation(op code.Opcode, left, right object.Object) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown operator: %c (%s %s)", op, left.Type(), right.Type())
	}
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return vm.push(&object.String{Value: leftVal + rightVal})
}

// executeComparison 执行比较操作
func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	leftType := left.Type()
	rightType := right.Type()
	if leftType == object.INTEGER && rightType == object.INTEGER {
		return vm.executeIntegerComparison(op, left, right)
	}
	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(left == right))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(left != right))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, leftType, rightType)
	}
}

// executeIntegerComparison 执行整数比较
func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	var result bool
	switch op {
	case code.OpEqual:
		result = leftVal == rightVal
	case code.OpNotEqual:
		result = leftVal != rightVal
	case code.OpGreaterThan:
		result = leftVal > rightVal
	default:
		return fmt.Errorf("unknown operator: %c", op)
	}
	return vm.push(nativeBoolToBooleanObject(result))

}

// nativeBoolToBooleanObject 将布尔值转换为布尔对象
func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return True
	}
	return False
}

// executeBangOperator 执行逻辑非操作
func (vm *VM) executeBangOperator() error {
	operand := vm.pop()
	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

// executeMinusOperator 执行负号操作
func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()
	if operand.Type() != object.INTEGER {
		return fmt.Errorf("unsupported type for negation: %s", operand.Type())
	}
	value := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -value})
}

// isTruthy 判断对象是否为真
func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}

// buildArray 从栈中构建一个数组对象
func (vm *VM) buildArray(startIndex, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)
	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}
	return &object.Array{Elements: elements}
}
