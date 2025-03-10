package compiler

import (
	"fmt"

	"monkey/ast"
	"monkey/code"
	"monkey/object"
)

// Compiler 编译器
type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

// New 创建编译器
func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

// Compile 编译
func (c *Compiler) Compile(node ast.Node) error {
	switch n := node.(type) {
	case *ast.Program:
		for _, s := range n.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(n.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)
	case *ast.InfixExpression:
		err := c.Compile(n.Left)
		if err != nil {
			return err
		}
		err = c.Compile(n.Right)
		if err != nil {
			return err
		}
		switch n.Operator {
		case "+":
			c.emit(code.OpAdd)
		default:
			return fmt.Errorf("unsupported operator %s", n.Operator)
		}
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: n.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	}
	return nil
}

// addConstant 添加常量
func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// emit 添加指令
func (c *Compiler) emit(op code.Opcode, operand ...int) int {
	return c.addInstruction(code.Make(op, operand...))
}

// addInstruction 添加指令
func (c *Compiler) addInstruction(ins []byte) int {
	posNewIns := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewIns
}

// Bytecode 产生字节码
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

// Bytecode 字节码
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
