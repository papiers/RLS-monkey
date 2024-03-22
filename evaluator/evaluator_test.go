package evaluator

import (
	"testing"

	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

func TestEvalIntegerExpression(t *testing.T) {

	tests := []struct {
		input string
		want  int64
	}{
		{"1", 1},
		{"2", 2},
		{"3", 3},
		{"-1", -1},
		{"-2", -2},
		{"-3", -3},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"1 / 2", 0},
		{"-5 / 4", -1},
		{"-5 / -4", 1},
		{"10 / 3", 3},
		{"3*3*3", 27},
		{"3*3/3", 3},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.want)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("obj is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("obj has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"1 == 1", true},
		{"1 != 1", false},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Fatalf("obj is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("obj has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input string
		want  interface{}
	}{
		{"if (true) { 10 }", int64(10)},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", int64(10)},
		{"if (1 < 2) { 10 }", int64(10)},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", int64(20)},
		{"if (1 < 2) { 10 } else { 20 }", int64(10)},
		{"if (1 < 2) { if(true){10} }", int64(10)},
		{"if (1 < 2) { if(false){10} }", nil},
		{"if (1 < 2) { if(true){10} else{20} }", int64(10)},
		{"if (1 < 2) { if(false){10} else{20} }", int64(20)},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.want.(int64)
		if ok {
			testIntegerObject(t, evaluated, integer)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != Null {
		t.Errorf("obj is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unsupported operator: -BOOLEAN"},
		{"true + false;", "unsupported operator: BOOLEAN + BOOLEAN"},
		{"true + 10;", "type mismatch: BOOLEAN + INTEGER"},
		{"false + 10;", "type mismatch: BOOLEAN + INTEGER"},
		{"5; true + 10;", "type mismatch: BOOLEAN + INTEGER"},
		{"5; 10 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"if (10 > 1) { if (10 > 2) { return true + false; } return 1; } 10", "unsupported operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
		{`"Hello" - "World"`, "unsupported operator: STRING - STRING"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Logf("obj is not Error. got=%T (%+v)", evaluated, evaluated)
			t.FailNow()
		}
		if errObj.Message != tt.expected {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("obj is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Logf("fn.Parameters does not have 1 parameters. got=%d", len(fn.Parameters))
		t.FailNow()
	}
	if fn.Parameters[0].String() != "x" {
		t.Errorf("parameter is not x. got=%q", fn.Parameters[0])
	}
	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Errorf("body is not %q. got=%q", expectedBody, fn.Body)
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(10);", 20},
		{"let add = fn(x, y) { x + y; }; add(5, 10);", 15},
		{"let add = fn(x, y) { x + y; }; add(5 + 10, add(5, 5) + 10);", 35},
		{"fn(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("obj is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("str is not Hello World!. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	expected := "Hello World!"
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("obj is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != expected {
		t.Errorf("str is not %q. got=%q", expected, str.Value)
	}
}
