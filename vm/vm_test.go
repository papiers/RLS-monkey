package vm

import (
	"fmt"
	"testing"

	"monkey/ast"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

type vmTestCase struct {
	input    string
	expected any
}

// testIntegerArithmetic 测试整数算术
func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 * (2 + 10)", 60},
		{"5 * (2 + 10) * 2 + 10", 130},
		{"5 * (2 + 10) * 2 + 10 - 5", 125},
		{"-5", -5},
		{"-10", -10},
		{"50 / 2 * 2 + 10 + -5", 55},
	}
	runVMTests(t, tests)
}

// TestBooleanExpressions 测试布尔表达式
func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 > 2", false},
		{"1 < 2", true},
		{"1 > 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true != false", true},
		{"false != true", true},
		{"true != false", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 == 1) == true", true},
		{"(1 == 1) == false", false},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	runVMTests(t, tests)
}

// runVMTests 运行虚拟机测试
func runVMTests(t *testing.T, tests []vmTestCase) {
	t.Helper()
	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}
		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}
		stackElem := vm.LastPoppedStackElem()
		testExpectedObject(t, tt.expected, stackElem)
	}
}

// testExpectedObject 测试期望的对象
func testExpectedObject(t *testing.T, expected any, actual object.Object) {
	t.Helper()
	switch exp := expected.(type) {
	case int:
		err := testIntegerObject(int64(exp), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(exp, actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	default:
		t.Errorf("type of expected value not handled. Got=%T", exp)
	}
}

// parses 解析输入的源代码
func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

// testIntegerObject 测试整数对象
func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer: %T", result)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value, got %d want %d", result.Value, expected)
	}
	return nil
}

// testBooleanObject 测试布尔对象
func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", result, result)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value, got %t want %t", result.Value, expected)
	}
	return nil
}
