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
		{"!(if (false) { 5; })", true},
	}
	runVMTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", Null},
		{"if (false) { 10 }", Null},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
	}
	runVMTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one + two", 3},
		{"let one = 1; let two = one + one; one + two", 3},
	}
	runVMTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}
	runVMTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1, 2, 3]", []int{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}
	runVMTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"{}", map[object.HashKey]int64{}},
		{
			"{1: 2, 2:3}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 1}).HashKey(): 2,
				(&object.Integer{Value: 2}).HashKey(): 3,
			},
		},
		{
			"{1+1: 2*2, 3+3:4*4}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 2}).HashKey(): 4,
				(&object.Integer{Value: 6}).HashKey(): 16,
			},
		},
	}
	runVMTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[1,2,3][1]", 2},
		{"[1,2,3][0+2]", 3},
		{"[[1,1,1]][0][0]", 1},
		{"[][0]", Null},
		{"[1,2,3][99]", Null},
		{"[1][-1]", Null},
		{"{1:1, 2:2}[1]", 1},
		{"{1:1, 2:2}[2]", 2},
		{"{1:1}[0]", Null},
		{"{}[0]", Null},
	}
	runVMTests(t, tests)
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let fivePlusTen = fn() {5+10;};
			fivePlusTen();
			`,
			expected: 15,
		},
		{
			input: `
			let one = fn() {1;};
			let two = fn() {2;};
			one() + two()
			`,
			expected: 3,
		},
		{
			input: `
			let a = fn() {1};
			let b = fn() {a() + 1};
			let c = fn() {b() + 1};
			c();
			`,
			expected: 3,
		},
	}
	runVMTests(t, tests)
}

func TestFunctionsWithReturnStatement(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let earlyExit = fn() {return 99; 100;};
			earlyExit();
			`,
			expected: 99,
		},
		{
			input: `
			let earlyExit = fn() {return 99; return 100;};
			earlyExit();
			`,
			expected: 99,
		},
	}
	runVMTests(t, tests)
}

func TestFunctionsWithoutReturnValue(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let noReturn = fn() {};
			noReturn();
			`,
			expected: Null,
		},
		{
			input: `
			let noReturn = fn() {};
			let noReturnTwo = fn() {noReturn();};
			noReturn();
			noReturnTwo();
			`,
			expected: Null,
		},
	}
	runVMTests(t, tests)
}

func TestFirstClassFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let returnsOne = fn() {1;};
			let returnsOneReturner = fn() {returnsOne;};
			returnsOneReturner()();
			`,
			expected: 1,
		},
		{
			input: `
			let returnsOneReturner = fn() {
				let returnsOne = fn() { 1; };
				returnsOne;
			};
			returnsOneReturner()();
			`,
			expected: 1,
		},
	}
	runVMTests(t, tests)
}

func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
		let one = fn() { let one = 1; one };
		one();
		`,
			expected: 1,
		},
		{
			input: `
		let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
		oneAndTwo();
		`,
			expected: 3,
		},
		{
			input: `
		let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
		let threeAndFour = fn() { let three = 3; let four = 4; three + four; };
		oneAndTwo() + threeAndFour();
		`,
			expected: 10,
		},
		{
			input: `
		let firstFoobar = fn() { let foobar = 50; foobar; };
		let secondFoobar = fn() { let foobar = 100; foobar; };
		firstFoobar() + secondFoobar();
		`,
			expected: 150,
		},
		{
			input: `
		let globalSeed = 50;
		let minusOne = fn() {
			let num = 1;
			globalSeed - num;
		}
		let minusTwo = fn() {
			let num = 2;
			globalSeed - num;
		}
		minusOne() + minusTwo();
		`,
			expected: 97,
		},
	}

	runVMTests(t, tests)
}

func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
		let identity = fn(a) { a; };
		identity(4);
		`,
			expected: 4,
		},
		{
			input: `
		let sum = fn(a, b) { a + b; };
		sum(1, 2);
		`,
			expected: 3,
		},
		{
			input: `
		let sum = fn(a, b) {
			let c = a + b;
			c;
		};
		sum(1, 2);
		`,
			expected: 3,
		},
		{
			input: `
		let sum = fn(a, b) {
			let c = a + b;
			c;
		};
		sum(1, 2) + sum(3, 4);`,
			expected: 10,
		},
		{
			input: `
		let sum = fn(a, b) {
			let c = a + b;
			c;
		};
		let outer = fn() {
			sum(1, 2) + sum(3, 4);
		};
		outer();
		`,
			expected: 10,
		},
		{
			input: `
			let globalNum = 10;
			let sum = fn(a,b) {
				let c = a + b;
				c + globalNum;
			};
			let outer = fn() {
				sum(1,2) + sum(3,4) + globalNum;
			};
			outer() + globalNum;
			`,
			expected: 50,
		},
	}

	runVMTests(t, tests)
}

func TestCallingFunctionsWithWrongArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `fn() {1;}(1);`,
			expected: `wrong number of arguments: want=0, got=1`,
		},
		{
			input:    `fn(a) {a;}();`,
			expected: `wrong number of arguments: want=1, got=0`,
		},
		{
			input:    `fn(a,b) {a+b;}(1);`,
			expected: `wrong number of arguments: want=2, got=1`,
		},
	}
	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}
		vm := New(comp.Bytecode())
		err = vm.Run()
		if err == nil {
			t.Fatalf("expected VM error but resulted in none.")
		}
		if err.Error() != tt.expected {
			t.Fatalf("wrong VM error: want=%q, got=%q", tt.expected, err)
		}
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []vmTestCase{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{
			`len(1)`,
			&object.Error{
				Message: "argument to `len` not supported, got INTEGER",
			},
		},
		{`len("one", "two")`,
			&object.Error{
				Message: "wrong number of arguments. got=2, want=1",
			},
		},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`puts("hello", "world!")`, Null},
		{`first([1, 2, 3])`, 1},
		{`first([])`, Null},
		{`first(1)`,
			&object.Error{
				Message: "argument to `first` must be ARRAY, got INTEGER",
			},
		},
		{`last([1, 2, 3])`, 3},
		{`last([])`, Null},
		{`last(1)`,
			&object.Error{
				Message: "argument to `last` must be ARRAY, got INTEGER",
			},
		},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest([])`, Null},
		{`push([], 1)`, []int{1}},
		{`push(1, 1)`,
			&object.Error{
				Message: "argument to `push` must be ARRAY, got INTEGER",
			},
		},
	}

	runVMTests(t, tests)
}

func TestClosures(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let newClosure = fn(a) {
				fn() {a;};
			};
			let closure = newClosure(99);
			closure();
			`,
			expected: 99,
		},
		{
			input: `
			let newAdder = fn(a,b) {
				fn(c) {a+b+c};
			};
			let adder = newAdder(1,2);
			adder(8);
			`,
			expected: 11,
		},
		{
			input: `
			let newAdder = fn(a,b) {
				let c = a+b;
				fn(d) {c+d};
			};
			let adder = newAdder(1,2);
			adder(8);
			`,
			expected: 11,
		},
		{
			input: `
		let newAdderOuter = fn(a, b) {
			let c = a + b;
			fn(d) {
				let e = d + c;
				fn(f) { e + f; };
			};
		};
		let newAdderInner = newAdderOuter(1, 2)
		let adder = newAdderInner(3);
		adder(8);
		`,
			expected: 14,
		},
		{
			input: `
		let a = 1;
		let newAdderOuter = fn(b) {
			fn(c) {
				fn(d) { a + b + c + d };
			};
		};
		let newAdderInner = newAdderOuter(2)
		let adder = newAdderInner(3);
		adder(8);
		`,
			expected: 14,
		},
		{
			input: `
		let newClosure = fn(a, b) {
			let one = fn() { a; };
			let two = fn() { b; };
			fn() { one() + two(); };
		};
		let closure = newClosure(9, 90);
		closure();
		`,
			expected: 99,
		},
	}
	runVMTests(t, tests)
}

func TestRecursiveFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let countDown = fn(x) {
				if (x==0) {
					return 0;
				} else {
					countDown(x-1);
				}
			};
			countDown(1);
			`,
			expected: 0,
		},
		{
			input: `
			let countDown = fn(x) {
				if (x == 0) {
					return 0
				} else {
					countDown(x-1);
				}
			};
			let wrapper = fn() {
				countDown(2);
			};
			wrapper();
			`,
			expected: 0,
		},
		{
			input: `
			let wrapper = fn() {
				let countDown = fn(x) {
					if (x == 0) {
						return 0
					} else {
						countDown(x-1);
					}
				};
				countDown(2);
			};
			wrapper();
			`,
			expected: 0,
		},
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
	case string:
		err := testStringObject(exp, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case []int:
		array, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("object not Array: %T (%+v)", actual, actual)
			return
		}
		if len(array.Elements) != len(exp) {
			t.Errorf("wrong num of elements. want=%d, got=%d", len(exp), len(array.Elements))
			return
		}
		for i, expectedElem := range exp {
			err := testIntegerObject(int64(expectedElem), array.Elements[i])
			if err != nil {
				t.Errorf("testIntergerObject failed: %s", err)
			}
		}
	case map[object.HashKey]int64:
		hash, ok := actual.(*object.Hash)
		if !ok {
			t.Errorf("object is not Hash. got=%T (%+v)", actual, actual)
			return
		}
		if len(hash.Pairs) != len(exp) {
			t.Errorf("hash has wrong num of pairs. got=%d want=%d", len(hash.Pairs), len(exp))
			return
		}
		for expectedKey, expectedValue := range exp {
			pair, ok := hash.Pairs[expectedKey]
			if !ok {
				t.Errorf("no pair for given key in hash")
				return
			}
			err := testIntegerObject(expectedValue, pair.Value)
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}
	case *object.Null:
		if actual != Null {
			t.Errorf("object is not NULL. got=%T (%+v)", actual, actual)
		}
	case *object.Error:
		errObj, ok := actual.(*object.Error)
		if !ok {
			t.Errorf("object is not Error. got=%T (%+v)", actual, actual)
			return
		}
		if len(errObj.Message) != len(exp.Message) {
			t.Errorf("object has wrong value. got=%q, want=%q", errObj.Message, exp.Message)
			return
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

// testStringObject 测试字符串对象
func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
	}
	return nil
}
