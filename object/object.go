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
	STRING       TypeObject = "STRING"
	BUILTIN      TypeObject = "BUILTIN"
	ARRAY        TypeObject = "ARRAY"
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

// 定义 Integer 对象实现 Object 接口
var _ Object = (*Integer)(nil)

// Type 返回对象类型
func (i *Integer) Type() TypeObject { return INTEGER }

// Inspect 返回对象字符串表示
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Boolean 布尔对象
type Boolean struct {
	Value bool
}

// 定义 Boolean 对象实现 Object 接口
var _ Object = (*Boolean)(nil)

// Type 返回对象类型
func (b *Boolean) Type() TypeObject { return BOOLEAN }

// Inspect 返回对象字符串表示
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// Null 空对象
type Null struct{}

// 定义 Null 对象实现 Object 接口
var _ Object = (*Null)(nil)

// Type 返回对象类型
func (*Null) Type() TypeObject { return NULL }

// Inspect 返回对象字符串表示
func (*Null) Inspect() string { return "null" }

// ReturnValue 返回对象
type ReturnValue struct {
	Value Object
}

// 定义 ReturnValue 对象实现 Object 接口
var _ Object = (*ReturnValue)(nil)

// Type 返回对象类型
func (rv *ReturnValue) Type() TypeObject { return RETURN_VALUE }

// Inspect 返回对象字符串表示
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Error 错误对象
type Error struct {
	Message string
}

// 定义 Error 对象实现 Object 接口
var _ Object = (*Error)(nil)

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

// 定义 Function 对象实现 Object 接口
var _ Object = (*Function)(nil)

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

// String 字符串对象
type String struct {
	Value string
}

// 定义 String 对象实现 Object 接口
var _ Object = (*String)(nil)

// Type 返回对象类型
func (s *String) Type() TypeObject { return STRING }

// Inspect 返回对象字符串表示
func (s *String) Inspect() string { return s.Value }

// BuiltinFunction 自定义函数
type BuiltinFunction func(args ...Object) Object

// Builtin 自定义函数对象
type Builtin struct {
	Fn BuiltinFunction // 自定义函数
}

// 定义 Builtin 对象实现 Object 接口
var _ Object = (*Builtin)(nil)

// Type 返回对象类型
func (b *Builtin) Type() TypeObject { return BUILTIN }

// Inspect 返回对象字符串表示
func (b *Builtin) Inspect() string { return "builtin function" }

// Array 数组对象
type Array struct {
	Elements []Object
}

// 定义 Array 对象实现 Object 接口
var _ Object = (*Array)(nil)

// Type 返回对象类型
func (a *Array) Type() TypeObject { return ARRAY }

// Inspect 返回对象字符串表示
func (a *Array) Inspect() string {
	var out strings.Builder
	elements := make([]string, len(a.Elements))
	for i, element := range a.Elements {
		elements[i] = element.Inspect()
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
