package parser

import (
	"fmt"
	"testing"

	"monkey/ast"
	"monkey/lexer"
)

func TestStatement(t *testing.T) {
	input := `
		let x1 = 5;
		let y = 10;
		let foobar = 838383;
		return 5 + y * 10;
		return 10 + x;
		return 99; 
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 6 {
		t.Fatalf("program.Statements does not contain 6 statements. got=%d\n", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		expectedVal        string
	}{
		{"x1", "5"},
		{"y", "10"},
		{"foobar", "838383"},
		{"", ""},
		{"", ""},
		{"", ""},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if tt.expectedIdentifier == "" {
			if !testReturnStatement(t, stmt) {
				return
			}
		} else {
			if !testLetStatement(t, stmt, tt.expectedIdentifier) {
				return
			}
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d\n", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", stmt)
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %q. got=%q", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %q. got=%q", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d\n", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", stmt)
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %q. got=%q", "5", literal.TokenLiteral())
	}

}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue any
	}{
		{input: `!5;`, operator: "!", integerValue: 5},
		{input: `-15;`, operator: "-", integerValue: 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d\n", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", stmt)
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %q. got=%q", tt.operator, exp.Operator)
		}
		if !testLiteralsExpression(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{input: "5 + 5;", leftValue: 5, operator: "+", rightValue: 5},
		{input: "5 - 5;", leftValue: 5, operator: "-", rightValue: 5},
		{input: "5 * 5;", leftValue: 5, operator: "*", rightValue: 5},
		{input: "5 / 5;", leftValue: 5, operator: "/", rightValue: 5},
		{input: "5 > 5;", leftValue: 5, operator: ">", rightValue: 5},
		{input: "5 < 5;", leftValue: 5, operator: "<", rightValue: 5},
		{input: "5 == 5;", leftValue: 5, operator: "==", rightValue: 5},
		{input: "5 != 5;", leftValue: 5, operator: "!=", rightValue: 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d\n", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", stmt)
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Expressions[0] is not ast.InfixExpression. got=%T", stmt.Expression)
		}
		if !testInfixExpression(t, exp, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorsPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b", "((-a) * b)",
		},
		{
			"!-a", "(!(-a))",
		},
		{
			"a + b + c", "((a + b) + c)",
		},
		{
			"a + b - c", "((a + b) - c)",
		},
		{
			"a * b * c", "((a * b) * c)",
		},
		{
			"a + b / c", "(a + (b / c))",
		},
		{
			"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))",
		},
		{
			"2 + 3 * 4 == 3 - 6 / 2", "((2 + (3 * 4)) == (3 - (6 / 2)))",
		},
		{
			"true == true", "(true == true)",
		},
		{
			"false == false", "(false == false)",
		},
		{
			"3 > 5 == false", "((3 > 5) == false)",
		},
		{
			"3 < 4 != 2 < 4", "((3 < 4) != (2 < 4))",
		},
		{
			"1 + 2 * 3 + 4", "((1 + (2 * 3)) + 4)",
		},
		{
			"(5 + 4) * 2", "((5 + 4) * 2)",
		},
		{
			"2 / (5 + 4)", "(2 / (5 + 4))",
		},
		{
			"-(5 + 4)", "(-(5 + 4))",
		},
		{
			"!(true == true)", "(!(true == true))",
		},
		{
			"a + add(b * c) + d", "((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, sub(6, 7))", "add(a, b, 1, (2 * 3), (4 + 5), sub(6, 7))",
		},
		{
			"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a*b[2], b[1], 2*3, 4+5, sub(6,7))", "add((a * (b[2])), (b[1]), (2 * 3), (4 + 5), sub(6, 7))",
		},
		{
			"[][a]", "([][a])",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got %d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. Got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. Got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence block statement is not ast.ExpressionStatement. Got=%T", exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v\n", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got %d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. Got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. Got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence block statement is not ast.ExpressionStatement. Got=%T", exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("alternative is not 1 statements. got=%d\n", len(exp.Alternative.Statements))
	}
	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative block statement is not ast.ExpressionStatement. Got=%T", exp.Alternative.Statements[0])
	}
	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d\n", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("First statement is not ast.ExpressionStatement. Got=%T", program.Statements[0])
	}
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("exp not *ast.FunctionLiteral. Got=%T", stmt.Expression)
	}
	if len(function.Parameters) != 2 {
		t.Errorf("function literal parameters wrong. got=%d\n", len(function.Parameters))
	}
	testLiteralsExpression(t, function.Parameters[0], "x")
	testLiteralsExpression(t, function.Parameters[1], "y")
	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body does not have enough statements. got=%d\n", len(function.Body.Statements))
	}
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. Got=%T", function.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{`fn() {};`, []string{}},
		{`fn(x) {};`, []string{"x"}},
		{`fn(x, y, z) {};`, []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)
		if len(function.Parameters) != len(tt.expected) {
			t.Errorf("length parameters wrong. got=%d, want=%d\n", len(function.Parameters), len(tt.expected))
		}
		for i, ident := range tt.expected {
			testLiteralsExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d\n", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("First statement is not ast.ExpressionStatement. Got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("exp not *ast.CallExpression. Got=%T", stmt.Expression)
	}
	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d\n", len(exp.Arguments))
	}
	testIntegerLiteral(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp is not *ast.StringLiteral. Got=%T", stmt.Expression)
	}
	if literal.Value != "hello world" {
		t.Errorf("literal.Value is %q, want %q", literal.Value, "hello world")
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp is not ast.ArrayLiteral. Got=%T", stmt.Expression)
	}
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) is not 3. Got=%d", len(array.Elements))
	}
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[123]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp is not ast.IndexExpression. Got=%T", stmt.Expression)
	}
	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}
	if !testIntegerLiteral(t, indexExp.Index, 123) {
		return
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hashLiteral, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. Got=%T", stmt.Expression)
	}
	if len(hashLiteral.Pairs) != 2 {
		t.Errorf("len(hashLiteral.Pairs) is not 2. Got=%d", len(hashLiteral.Pairs))
	}
	expected := map[string]int64{
		"one": 1,
		"two": 2,
	}
	for keyNode, valueNode := range hashLiteral.Pairs {
		key, ok := keyNode.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not *ast.StringLiteral. Got=%T", keyNode)
		}
		expectedKey := expected[key.Value]
		testIntegerLiteral(t, valueNode, expectedKey)
	}
}

func TestParsingEmptyHastLiterals(t *testing.T) {
	input := "{}"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hashLiteral, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. Got=%T", stmt.Expression)
	}
	if len(hashLiteral.Pairs) != 0 {
		t.Errorf("len(hashLiteral.Pairs) is not 0. Got=%d", len(hashLiteral.Pairs))
	}
}

func TestParsingHashLiteralsExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 2 * 3}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hashLiteral, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. Got=%T", stmt.Expression)
	}
	if len(hashLiteral.Pairs) != 3 {
		t.Errorf("len(hashLiteral.Pairs) is not 3. Got=%d", len(hashLiteral.Pairs))
	}
	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 2, "*", 3)
		},
	}
	for keyNode, valueNode := range hashLiteral.Pairs {
		key, ok := keyNode.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not *ast.StringLiteral. Got=%T", keyNode)
			continue
		}
		testFunc, ok := tests[key.Value]
		if !ok {
			t.Errorf("No test function for %q key.", key.Value)
			continue
		}
		testFunc(valueNode)
	}
}

func TestFunctionLiteralWithName(t *testing.T) {
	input := `let myFunction = fn() {};`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.LetStatement. got=%T",
			program.Statements[0])
	}
	function, ok := stmt.Value.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not ast.FunctionLiteral. got=%T",
			stmt.Value)
	}
	if function.Name != "myFunction" {
		t.Fatalf("function literal name wrong. want 'myFunction', got=%q\n",
			function.Name)
	}
}

// testLetStatement 测试解析let表达式
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not let statement. got=%q", s)
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not %q. got=%q", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral not %q. got=%q", name, letStmt.Name.TokenLiteral())
	}

	return true
}

// testReturnStatement 测试解析return表达式
func testReturnStatement(t *testing.T, s ast.Statement) bool {
	if s.TokenLiteral() != "return" {
		t.Errorf("s.TokenLiteral not return statement. got=%q", s)
		return false
	}
	_, ok := s.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("s not *ast.ReturnStatement. got=%T", s)
		return false
	}

	return true
}

// checkParserErrors 解析器错误检查
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

// testIntegerLiteral 测试解析整数表达式
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	intLit, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if intLit.Value != value {
		t.Errorf("il.Value not %d. got=%d", value, intLit.Value)
		return false
	}
	if intLit.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("intLit.TokenLiteral not %q. got=%q", value, intLit.TokenLiteral())
		return false
	}
	return true
}

// testIdentifier 测试解析标识符表达式
func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %q. got=%q", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %q. got=%q", value, ident.TokenLiteral())
		return false
	}
	return true
}

// testBooleanLiteral 测试解析布尔表达式
func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())
		return false
	}
	return true
}

// testLiteralsExpression 测试解析字面量表达式
func testLiteralsExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", expected)
	return false
}

// testInfixExpression 测试解析中缀表达式
func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%+v)", exp, exp)
		return false
	}
	if !testLiteralsExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not %q. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralsExpression(t, opExp.Right, right) {
		return false
	}
	return true
}
