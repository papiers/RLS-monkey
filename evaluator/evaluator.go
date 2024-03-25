package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
	True = &object.Boolean{
		Value: true,
	}
	False = &object.Boolean{
		Value: false,
	}
	Null = &object.Null{}
)

// Eval 执行表达式
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	default:
		return &object.Error{Message: "unknown node type for eval"}
	}
	return nil
}

// evalPrefixExpression 执行前缀表达式
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return &object.Error{
			Message: "unsupported operator: " + operator + string(right.Type()),
		}
	}
}

// evalProgram 执行语句列表
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

// evalBlockStatement 执行块语句
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			switch rt {
			case object.RETURN_VALUE:
				return result
			case object.ERROR:
				return result
			}
		}
	}
	return result
}

// nativeBoolToBooleanObject 将布尔值转换为 Monkey 对象
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}
	return False
}

// evalBangOperatorExpression 执行前缀表达式 !
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case True:
		return False
	case False:
		return True
	case Null:
		return True
	default:
		return False
	}
}

// evalMinusPrefixOperatorExpression 执行前缀表达式 -
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if integer, ok := right.(*object.Integer); ok && right.Type() == object.INTEGER {
		return &object.Integer{Value: -integer.Value}
	}
	return &object.Error{
		Message: "unsupported operator: -" + string(right.Type()),
	}
}

// evalInfixExpression 执行中缀表达式
func evalInfixExpression(operator string, left, right object.Object) object.Object {

	if left.Type() == object.INTEGER && right.Type() == object.INTEGER {
		l, okLeft := left.(*object.Integer)
		r, okRight := right.(*object.Integer)
		if okLeft && okRight {
			return evalIntegerInfixExpression(operator, l, r)
		}
	}
	if left.Type() == object.STRING && right.Type() == object.STRING {
		l, okLeft := left.(*object.String)
		r, okRight := right.(*object.String)
		if okLeft && okRight {
			return evalStringInfixExpression(operator, l, r)
		}
	}
	if operator == "==" {
		return nativeBoolToBooleanObject(left == right)
	} else if operator == "!=" {
		return nativeBoolToBooleanObject(left != right)
	}
	if left.Type() != right.Type() {
		return &object.Error{Message: "type mismatch: " + string(left.Type()) + " " + operator + " " + string(right.Type())}
	}
	return &object.Error{Message: "unsupported operator: " + string(left.Type()) + " " + operator + " " + string(right.Type())}
}

// evalIntegerInfixExpression 执行中缀表达式，整数类型
func evalIntegerInfixExpression(operator string, left, right *object.Integer) object.Object {
	switch operator {
	case "+":
		return &object.Integer{Value: left.Value + right.Value}
	case "-":
		return &object.Integer{Value: left.Value - right.Value}
	case "*":
		return &object.Integer{Value: left.Value * right.Value}
	case "/":
		return &object.Integer{Value: left.Value / right.Value}
	case "<":
		return nativeBoolToBooleanObject(left.Value < right.Value)
	case ">":
		return nativeBoolToBooleanObject(left.Value > right.Value)
	case "==":
		return nativeBoolToBooleanObject(left.Value == right.Value)
	case "!=":
		return nativeBoolToBooleanObject(left.Value != right.Value)
	}
	return &object.Error{Message: "unsupported operator: " + string(left.Type()) + " " + operator + " " + string(right.Type())}
}

// evalStringInfixExpression
func evalStringInfixExpression(operator string, left, right *object.String) object.Object {
	switch operator {
	case "+":
		return &object.String{Value: left.Value + right.Value}
	}
	return &object.Error{Message: "unsupported operator: " + string(left.Type()) + " " + operator + " " + string(right.Type())}
}

// evalIfExpression 计算if表达式
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}
	return Null
}

// isTruthy 判断对象是否为真
func isTruthy(obj object.Object) bool {
	if obj == Null {
		return false
	}
	if obj == True {
		return true
	}
	if obj == False {
		return false
	}
	return true
}

// isError 判断对象是否为错误
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}

// evalIdentifier 计算标识符
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return &object.Error{Message: "identifier not found: " + node.Value}
}

// evalExpressions 计算表达式列表
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

// applyFunction 计算函数调用
func applyFunction(fn object.Object, args []object.Object) object.Object {
	if fun, ok := fn.(*object.Function); ok {
		extendedEnv := extendFunctionEnv(fun, args)
		evaluated := Eval(fun.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	}

	if builtin, ok := fn.(*object.Builtin); ok {
		return builtin.Fn(args...)
	}

	return &object.Error{Message: "not a function"}
}

// extendFunctionEnv 扩展函数环境
func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

// unwrapReturnValue 解包返回值
func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// evalIndexExpression 计算索引表达式
func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		l, okL := left.(*object.Array)
		i, okI := index.(*object.Integer)
		if okL && okI {
			return evalArrayIndexExpression(l, i)
		}
	case left.Type() == object.HASH:
		l, okL := left.(*object.Hash)
		if okL {
			return evalHashIndexExpression(l, index)
		}
	}
	return &object.Error{Message: "index operator not supported"}
}

// evalArrayIndexExpression 计算数组索引表达式
func evalArrayIndexExpression(arr *object.Array, index *object.Integer) object.Object {
	i := int(index.Value)
	if i < 0 || i > len(arr.Elements)-1 {
		return Null
	}
	return arr.Elements[i]
}

// evalHashLiteral 计算哈希字面量
func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return &object.Error{Message: "unusable as hash key"}
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}
		pairs[hashKey.HashKey()] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}

// evalHashIndexExpression 计算哈希索引表达式
func evalHashIndexExpression(hash *object.Hash, index object.Object) object.Object {
	key, ok := index.(object.Hashable)
	if !ok {
		return &object.Error{Message: "unusable as hash key: " + string(index.Type())}
	}
	pair, ok := hash.Pairs[key.HashKey()]
	if !ok {
		return Null
	}
	return pair.Value
}
