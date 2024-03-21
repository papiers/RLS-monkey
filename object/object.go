package object

import (
	"fmt"
	"strings"

	"monkey/ast"
)

const (
	INTEGER      TypeObject = "INTEGER"
	BOOLEAN      TypeObject = "BOOLEAN"
	NULL         TypeObject = "NULL"
	RETURN_VALUE TypeObject = "RETURN_VALUE"
	ERROR        TypeObject = "ERROR"
	FUNCTION     TypeObject = "FUNCTION"
)

// TypeObject 对象类型
type TypeObject string

// Object 对象接口
type Object interface {
	Type() TypeObject
	Inspect() string
}

// Integer 整数对象
type Integer struct {
	Value int64
}

// Type 返回对象类型
func (i *Integer) Type() TypeObject { return INTEGER }

// Inspect 返回对象字符串表示
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Boolean 布尔对象
type Boolean struct {
	Value bool
}

// Type 返回对象类型
func (b *Boolean) Type() TypeObject { return BOOLEAN }

// Inspect 返回对象字符串表示
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// Null 空对象
type Null struct{}

// Type 返回对象类型
func (*Null) Type() TypeObject { return NULL }

// Inspect 返回对象字符串表示
func (*Null) Inspect() string { return "null" }

// ReturnValue 返回对象
type ReturnValue struct {
	Value Object
}

// Type 返回对象类型
func (rv *ReturnValue) Type() TypeObject { return RETURN_VALUE }

// Inspect 返回对象字符串表示
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Error 错误对象
type Error struct {
	Message string
}

// Type 返回对象类型
func (e *Error) Type() TypeObject { return ERROR }

// Inspect 返回对象字符串表示
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

// Function 函数对象
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Type 返回对象类型
func (f *Function) Type() TypeObject { return FUNCTION }

// Inspect 返回对象字符串表示
func (f *Function) Inspect() string {
	var out strings.Builder
	params := make([]string, len(f.Parameters))
	for i, param := range f.Parameters {
		params[i] = param.String()
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}
