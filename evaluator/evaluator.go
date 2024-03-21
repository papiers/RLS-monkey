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
	l, okLeft := left.(*object.Integer)
	r, okRight := right.(*object.Integer)
	if okLeft && okRight && left.Type() == object.INTEGER && right.Type() == object.INTEGER {
		return evalIntegerInfixExpression(operator, l, r)
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
	fun, ok := fn.(*object.Function)
	if !ok {
		return &object.Error{Message: "not a function"}
	}
	extendedEnv := extendFunctionEnv(fun, args)
	evaluated := Eval(fun.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
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
