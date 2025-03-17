package compiler

import (
	"fmt"
	"sort"

	"monkey/ast"
	"monkey/code"
	"monkey/object"
)

// EmittedInstruction 存储指令和位置
type EmittedInstruction struct {
	OpCode   code.Opcode
	Position int
}

// CompilationScope 编译作用域
type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

// Compiler 编译器
type Compiler struct {
	constants   []object.Object
	symbolTable *SymbolTable
	scopes      []CompilationScope
	scopeIndex  int
}

// New 创建编译器
func New() *Compiler {
	symbolTable := NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}
	return &Compiler{
		constants:   []object.Object{},
		symbolTable: symbolTable,
		scopes: []CompilationScope{
			{
				instructions:        code.Instructions{},
				lastInstruction:     EmittedInstruction{},
				previousInstruction: EmittedInstruction{},
			},
		},
		scopeIndex: 0,
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
	case *ast.StringLiteral:
		str := &object.String{Value: n.Value}
		c.emit(code.OpConstant, c.addConstant(str))
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
		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}
		jumpPos := c.emit(code.OpJump, 9999)

		// 回填else语句开始位置
		afterConsequencePos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if n.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			err = c.Compile(n.Alternative)
			if err != nil {
				return err
			}
			if c.lastInstructionIs(code.OpPop) {
				c.removeLastPop()
			}
		}
		// 回填else后语句开始位置
		afterAlternativePos := len(c.currentInstructions())
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
		if symbol.Scope == GlobalScope {
			c.emit(code.OpSetGlobal, symbol.Index)
		} else {
			c.emit(code.OpSetLocal, symbol.Index)
		}
	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(n.Value)
		if !ok {
			return fmt.Errorf("identifier not found: %s", n.Value)
		}
		c.loadSymbol(symbol)
	case *ast.ArrayLiteral:
		for _, e := range n.Elements {
			err := c.Compile(e)
			if err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(n.Elements))
	case *ast.HashLiteral:
		var keys []ast.Expression
		for k := range n.Pairs {
			keys = append(keys, k)
		}
		// 对键进行排序，以便在哈希表中保持一致的顺序
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})
		for _, v := range keys {
			err := c.Compile(v)
			if err != nil {
				return err
			}
			err = c.Compile(n.Pairs[v])
			if err != nil {
				return err
			}
		}
		c.emit(code.OpHash, len(n.Pairs)*2)
	case *ast.IndexExpression:
		err := c.Compile(n.Left)
		if err != nil {
			return err
		}
		err = c.Compile(n.Index)
		if err != nil {
			return err
		}
		c.emit(code.OpIndex)
	case *ast.FunctionLiteral:
		c.enterScope()
		for _, v := range n.Parameters {
			c.symbolTable.Define(v.Value)
		}
		err := c.Compile(n.Body)
		if err != nil {
			return err
		}
		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastPopWithReturn()
		}
		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}

		numLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()
		compiledFn := &object.CompiledFunction{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(n.Parameters),
		}
		c.emit(code.OpClosure, c.addConstant(compiledFn), 0)
	case *ast.ReturnStatement:
		err := c.Compile(n.ReturnValue)
		if err != nil {
			return err
		}
		c.emit(code.OpReturnValue)
	case *ast.CallExpression:
		err := c.Compile(n.Function)
		if err != nil {
			return err
		}
		for _, v := range n.Arguments {
			err = c.Compile(v)
			if err != nil {
				return err
			}
		}
		c.emit(code.OpCall, len(n.Arguments))
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
	posNewIns := len(c.currentInstructions())
	c.scopes[c.scopeIndex].instructions = append(c.currentInstructions(), ins...)
	return posNewIns
}

// setLastInstruction 设置最后一条指令
func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	c.scopes[c.scopeIndex].previousInstruction = c.scopes[c.scopeIndex].lastInstruction
	c.scopes[c.scopeIndex].lastInstruction = EmittedInstruction{op, pos}
}

// lastInstructionIsPop 最后一条指令是否为
func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}
	return c.scopes[c.scopeIndex].lastInstruction.OpCode == op
}

// removeLastPop 移除最后一条Pop
func (c *Compiler) removeLastPop() {
	c.scopes[c.scopeIndex].instructions = c.currentInstructions()[:c.scopes[c.scopeIndex].lastInstruction.Position]
	c.scopes[c.scopeIndex].lastInstruction = c.scopes[c.scopeIndex].previousInstruction
}

// replaceInstruction 替换指令
func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i, b := range newInstruction {
		c.scopes[c.scopeIndex].instructions[pos+i] = b
	}
}

// changeOperand 替换操作数
func (c *Compiler) changeOperand(pos int, operand int) {
	op := code.Opcode(c.currentInstructions()[pos])
	c.replaceInstruction(pos, code.Make(op, operand))
}

// currentInstructions 当前指令
func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

// enterScope 进入作用域
func (c *Compiler) enterScope() {
	c.scopes = append(c.scopes, CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	})
	c.scopeIndex++
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

// leaveScope 离开作用域
func (c *Compiler) leaveScope() code.Instructions {
	ins := c.currentInstructions()
	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--
	c.symbolTable = c.symbolTable.Outer
	return ins
}

// replaceLastPopWithReturn 替换最后一条Pop为Return
func (c *Compiler) replaceLastPopWithReturn() {
	lastIns := c.scopes[c.scopeIndex].lastInstruction
	c.replaceInstruction(lastIns.Position, code.Make(code.OpReturnValue))
	c.scopes[c.scopeIndex].lastInstruction.OpCode = code.OpReturnValue
}

// loadSymbol 加载符号
func (c *Compiler) loadSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(code.OpGetGlobal, s.Index)
	case LocalScope:
		c.emit(code.OpGetLocal, s.Index)
	case BuiltinScope:
		c.emit(code.OpGetBuiltin, s.Index)
	}
}

// Bytecode 产生字节码
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

// Bytecode 字节码
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
