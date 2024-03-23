package ast

import (
	"bytes"
	"strings"

	"monkey/token"
)

// Node 定义节点类型
type Node interface {
	TokenLiteral() string // 返回节点的token值
	String() string       // 返回节点的字符串
}

// Statement 定义语句节点类型
type Statement interface {
	Node            // 语句节点为节点
	statementNode() // 语句节点为语句
}

// Expression 定义表达式节点类型
type Expression interface {
	Node             // 表达式节点为节点
	expressionNode() // 表达式节点为表达式
}

// Program 定义程序节点
type Program struct {
	Statements []Statement // 程序节点中的语句
}

// 定义程序节点为节点
var _ Node = (*Program)(nil)

// TokenLiteral 返回程序节点的token值
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// String 返回程序节点的字符串
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Identifier 定义标识符节点
type Identifier struct {
	Token token.Token // 标识符的token
	Value string      // 标识符的值
}

// 定义标识符节点为表达式
var _ Expression = (*Identifier)(nil)

// expressionNode 标识符节点为表达式
func (i *Identifier) expressionNode() {}

// TokenLiteral 返回标识符的token值
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// String 返回标识符的字符串
func (i *Identifier) String() string {
	return i.Value
}

// LetStatement 定义let语句节点
type LetStatement struct {
	Token token.Token // let关键字token
	Name  *Identifier // let语句的标识符
	Value Expression  // let语句的值表达式
}

// 定义let语句节点为语句
var _ Statement = (*LetStatement)(nil)

// statementNode 标识let语句节点为语句
func (l *LetStatement) statementNode() {}

// TokenLiteral 返回let语句的token值
func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

// String 返回let语句的字符串
func (l *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String())
	out.WriteString(" = ")
	if l.Value != nil {
		out.WriteString(l.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement 定义return语句节点
type ReturnStatement struct {
	Token       token.Token // return关键字token
	ReturnValue Expression  // 返回值表达式
}

// 定义return语句节点为语句
var _ Statement = (*ReturnStatement)(nil)

// statementNode 标识return语句节点为语句
func (r ReturnStatement) statementNode() {
}

// TokenLiteral 返回return语句的token值
func (r ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

// String 返回return语句的字符串
func (r ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(r.TokenLiteral() + " ")
	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement 定义表达式语句节点
type ExpressionStatement struct {
	Token      token.Token // 表达式token
	Expression Expression  // 表达式
}

// 定义表达式语句节点为语句
var _ Statement = (*ExpressionStatement)(nil)

// statementNode 标识表达式语句节点为语句
func (e *ExpressionStatement) statementNode() {}

// TokenLiteral 返回表达式语句的token值
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

// String 返回表达式语句的字符串
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

// IntegerLiteral 定义整数节点
type IntegerLiteral struct {
	Token token.Token // 整数token
	Value int64       // 整数值
}

// 定义整数节点为表达式
var _ Expression = (*IntegerLiteral)(nil)

// expressionNode 标识整数节点为表达式
func (i *IntegerLiteral) expressionNode() {}

// TokenLiteral 返回整数节点的token值
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

// String 返回整数节点的字符串
func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

// PrefixExpression 定义前缀表达式节点
type PrefixExpression struct {
	Token    token.Token // 前缀表达式token
	Operator string      // 前缀运算符
	Right    Expression  // 右值表达式
}

// 定义前缀表达式节点为表达式
var _ Expression = (*PrefixExpression)(nil)

// expressionNode 标识前缀表达式节点为表达式
func (p *PrefixExpression) expressionNode() {}

// TokenLiteral 返回前缀表达式的token值
func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

// String 返回前缀表达式的字符串
func (p *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression 定义中缀表达式节点
type InfixExpression struct {
	Token    token.Token // 中缀表达式token
	Left     Expression  // 左值表达式
	Operator string      // 中缀运算符
	Right    Expression  // 右值表达式
}

// 定义中缀表达式节点为表达式
var _ Expression = (*InfixExpression)(nil)

// expressionNode 标识中缀表达式节点为表达式
func (i *InfixExpression) expressionNode() {}

// TokenLiteral 返回中缀表达式的token值
func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

// String 返回中缀表达式的字符串
func (i *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}

// Boolean 定义布尔节点
type Boolean struct {
	Token token.Token // 布尔token
	Value bool        // 布尔值
}

// 定义布尔节点为表达式
var _ Expression = (*Boolean)(nil)

// expressionNode 标识布尔节点为表达式
func (b *Boolean) expressionNode() {}

// TokenLiteral 返回布尔节点的token值
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

// String 返回布尔节点的字符串
func (b *Boolean) String() string {
	return b.Token.Literal
}

// BlockStatement 定义块语句节点
type BlockStatement struct {
	Token      token.Token // 块语句token
	Statements []Statement // 块语句节点列表
}

// 定义块语句节点为语句
var _ Statement = (*BlockStatement)(nil)

// statementNode 标识块语句节点为语句
func (b *BlockStatement) statementNode() {}

// TokenLiteral 返回块语句节点的token值
func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}

// String 返回块语句节点的字符串
func (b *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range b.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// IfExpression 定义if表达式节点
type IfExpression struct {
	Token       token.Token // if表达式token
	Condition   Expression  // 条件表达式
	Consequence *BlockStatement
	Alternative *BlockStatement
}

// 定义if表达式节点为表达式
var _ Expression = (*IfExpression)(nil)

// expressionNode 标识if表达式节点为表达式
func (i *IfExpression) expressionNode() {}

// TokenLiteral 返回if表达式的token值
func (i *IfExpression) TokenLiteral() string {
	return i.Token.Literal
}

// String 返回if表达式的字符串
func (i *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())
	if i.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(i.Alternative.String())
	}
	return out.String()
}

// FunctionLiteral 定义函数节点
type FunctionLiteral struct {
	Token      token.Token     // 函数token
	Parameters []*Identifier   // 函数参数列表
	Body       *BlockStatement // 函数体
}

// 定义函数节点为表达式
var _ Expression = (*FunctionLiteral)(nil)

// expressionNode 标识函数节点为表达式
func (f *FunctionLiteral) expressionNode() {}

// TokenLiteral 返回函数的token值
func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}

// String 返回函数的字符串
func (f *FunctionLiteral) String() string {
	var out bytes.Buffer
	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(f.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())
	return out.String()
}

// CallExpression 定义函数调用节点
type CallExpression struct {
	Token     token.Token  // 函数调用token
	Function  Expression   // 函数节点
	Arguments []Expression // 函数参数列表
}

// 定义函数调用节点为表达式
var _ Expression = (*CallExpression)(nil)

// expressionNode 标识函数调用节点为表达式
func (c *CallExpression) expressionNode() {}

// TokenLiteral 返回函数的token值
func (c *CallExpression) TokenLiteral() string {
	return c.Token.Literal
}

// String 返回函数的字符串
func (c *CallExpression) String() string {
	var out bytes.Buffer
	var args []string
	for _, a := range c.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// StringLiteral 定义字符串节点
type StringLiteral struct {
	Token token.Token // 字符串token
	Value string      // 字符串值
}

// 定义字符串节点为表达式
var _ Expression = (*StringLiteral)(nil)

// expressionNode 标识字符串节点为表达式
func (s *StringLiteral) expressionNode() {}

// TokenLiteral 返回字符串节点的token值
func (s *StringLiteral) TokenLiteral() string {
	return s.Token.Literal
}

// String 返回字符串节点的字符串
func (s *StringLiteral) String() string {
	return s.Token.Literal
}

// ArrayLiteral 定义数组节点
type ArrayLiteral struct {
	Token    token.Token  // 数组token
	Elements []Expression // 数组元素列表
}

// 定义数组节点为表达式
var _ Expression = (*ArrayLiteral)(nil)

// expressionNode 标识数组节点为表达式
func (a *ArrayLiteral) expressionNode() {}

// TokenLiteral 返回数组节点的token值
func (a *ArrayLiteral) TokenLiteral() string {
	return a.Token.Literal
}

// String 返回数组节点的字符串
func (a *ArrayLiteral) String() string {
	var out bytes.Buffer
	var elements []string
	for _, e := range a.Elements {
		elements = append(elements, e.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// IndexExpression 定义数组索引节点
type IndexExpression struct {
	Token token.Token // 索引token
	Left  Expression  // 数组节点
	Index Expression  // 索引节点
}

// 定义数组索引节点为表达式
var _ Expression = (*IndexExpression)(nil)

// expressionNode 标识数组索引节点为表达式
func (i *IndexExpression) expressionNode() {}

// TokenLiteral 返回数组索引节点的token值
func (i *IndexExpression) TokenLiteral() string {
	return i.Token.Literal
}

// String 返回数组索引节点的字符串
func (i *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString("[")
	out.WriteString(i.Index.String())
	out.WriteString("]")
	out.WriteString(")")
	return out.String()
}
