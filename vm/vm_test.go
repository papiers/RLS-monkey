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
