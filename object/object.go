package object

import (
	"fmt"
	"hash/fnv"
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
	HASH         TypeObject = "HASH"
)

// TypeObject 对象类型
type TypeObject string

// Object 对象接口
type Object interface {
	Type() TypeObject // 返回对象类型
	Inspect() string  // 返回对象字符串表示
}

// Hashable 哈希键接口
type Hashable interface {
	HashKey() HashKey
}

// Integer 整数对象
type Integer struct {
	Value int64 // 整数值
}

// 定义 Integer 对象实现 Object 接口
var _ Object = (*Integer)(nil)

// 定义 Integer 对象实现 Hashable 接口
var _ Hashable = (*Integer)(nil)

// Type 返回对象类型
func (i *Integer) Type() TypeObject { return INTEGER }

// Inspect 返回对象字符串表示
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// Boolean 布尔对象
type Boolean struct {
	Value bool // 布尔值
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
	Value Object // 返回值
}

// 定义 ReturnValue 对象实现 Object 接口
var _ Object = (*ReturnValue)(nil)

// Type 返回对象类型
func (rv *ReturnValue) Type() TypeObject { return RETURN_VALUE }

// Inspect 返回对象字符串表示
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Error 错误对象
type Error struct {
	Message string // 错误信息
}

// 定义 Error 对象实现 Object 接口
var _ Object = (*Error)(nil)

// Type 返回对象类型
func (e *Error) Type() TypeObject { return ERROR }

// Inspect 返回对象字符串表示
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

// Function 函数对象
type Function struct {
	Parameters []*ast.Identifier   // 参数列表
	Body       *ast.BlockStatement // 函数体
	Env        *Environment        // 函数执行环境
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
	Value string // 字符串值
}

// 定义 String 对象实现 Object 接口
var _ Object = (*String)(nil)

// Type 返回对象类型
func (s *String) Type() TypeObject { return STRING }

// Inspect 返回对象字符串表示
func (s *String) Inspect() string { return s.Value }

// HashKey 哈希键对象
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

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
	Elements []Object // 数组元素
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

// HashKey 哈希键对象
type HashKey struct {
	Type  TypeObject // 哈希键类型
	Value uint64     // 哈希键值
}

// HashPair 哈希键值对
type HashPair struct {
	Key   Object // 哈希键
	Value Object // 哈希值
}

// Hash 哈希对象
type Hash struct {
	Pairs map[HashKey]HashPair
}

// 定义 Hash 对象实现 Object 接口
var _ Object = (*Hash)(nil)

// Type 返回对象类型
func (h *Hash) Type() TypeObject { return HASH }

// Inspect 返回对象字符串表示
func (h *Hash) Inspect() string {
	var out strings.Builder
	pairs := make([]string, len(h.Pairs))
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
