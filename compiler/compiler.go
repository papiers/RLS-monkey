package compiler

import (
	"fmt"

	"monkey/ast"
	"monkey/code"
	"monkey/object"
)

// EmittedInstruction 存储指令和位置
type EmittedInstruction struct {
	OpCode   code.Opcode
	Position int
}

// Compiler 编译器
type Compiler struct {
	instructions        code.Instructions
	constants           []object.Object
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
	symbolTable         *SymbolTable
}

// New 创建编译器
func New() *Compiler {
	return &Compiler{
		instructions:        code.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
		symbolTable:         NewSymbolTable(),
	}
}

// NewWithState 创建编译器携带state
func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = s
	compiler.constants = constants
	return compiler
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
	case *ast.PrefixExpression:
		err := c.Compile(n.Right)
		if err != nil {
			return err
		}
		switch n.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unsupported prefix operator %s", n.Operator)
		}
	case *ast.InfixExpression:
		if n.Operator == "<" {
			err := c.Compile(n.Right)
			if err != nil {
				return err
			}
			err = c.Compile(n.Left)
			if err != nil {
				return err
			}
			c.emit(code.OpGreaterThan)
			return nil
		}
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
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unsupported operator %s", n.Operator)
		}
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: n.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	case *ast.Boolean:
		if n.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *ast.IfExpression:
		err := c.Compile(n.Condition)
		if err != nil {
			return err
		}
		// 记录发出虚假跳转指令的位置
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)
		err = c.Compile(n.Consequence)
		if err != nil {
			return err
		}
		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}
		jumpPos := c.emit(code.OpJump, 9999)

		// 回填else语句开始位置
		afterConsequencePos := len(c.instructions)
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if n.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			err = c.Compile(n.Alternative)
			if err != nil {
				return err
			}
			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}
		}
		// 回填else后语句开始位置
		afterAlternativePos := len(c.instructions)
		c.changeOperand(jumpPos, afterAlternativePos)

	case *ast.BlockStatement:
		for _, s := range n.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.LetStatement:
		err := c.Compile(n.Value)
		if err != nil {
			return err
		}
		symbol := c.symbolTable.Define(n.Name.Value)
		c.emit(code.OpSetGlobal, symbol.Index)
	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(n.Value)
		if !ok {
			return fmt.Errorf("identifier not found: %s", n.Value)
		}
		c.emit(code.OpGetGlobal, symbol.Index)
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
	ins := code.Make(op, operand...)
	pos := c.addInstruction(ins)
	c.setLastInstruction(op, pos)
	return pos
}

// addInstruction 添加指令
func (c *Compiler) addInstruction(ins []byte) int {
	posNewIns := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewIns
}

// setLastInstruction 设置最后一条指令
func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	c.previousInstruction = c.lastInstruction
	c.lastInstruction = EmittedInstruction{op, pos}
}

// lastInstructionIsPop 最后一条指令是否为Pop
func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.OpCode == code.OpPop
}

// removeLastPop 移除最后一条Pop
func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

// replaceInstruction 替换指令
func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i, b := range newInstruction {
		c.instructions[pos+i] = b
	}
}

// changeOperand 替换操作数
func (c *Compiler) changeOperand(pos int, operand int) {
	op := code.Opcode(c.instructions[pos])
	c.replaceInstruction(pos, code.Make(op, operand))
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
