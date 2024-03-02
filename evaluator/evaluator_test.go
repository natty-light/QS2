package evaluator

import (
	"QuonkScript/lexer"
	"QuonkScript/object"
	"QuonkScript/parser"
	"fmt"
	"testing"
)

func TestEvalIntegerExpr(t *testing.T) {
	tests := []struct {
		source   string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.source)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpr(t *testing.T) {
	tests := []struct {
		source   string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"true && false", false},
		{"true || false", true},
		{"(1 > 2) || (3 + 2 ==5)", true},
		{"3 > 2 && 4 > 3", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.source)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		source   string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.source)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
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
		{"if (10 > 1) { if (10 > 1) { return 10 } return 1 }", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
		expectedLine    int
	}{
		{
			"5 + true;",
			"type mismatch: Integer + Boolean",
			1,
		},
		{
			"5 + true; 5;",
			"type mismatch: Integer + Boolean",
			1,
		},
		{
			"-true",
			"unknown operation - for type Boolean",
			1,
		},
		{
			"true + false;",
			"unknown operator: Boolean + Boolean",
			1,
		},
		{
			"5; true + false; 5",
			"unknown operator: Boolean + Boolean",
			1,
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: Boolean + Boolean",
			1,
		},
		{
			` if (10 > 1) {
		  		if (10 > 1) {
					return true + false;
				}
				return 1;
			}`,
			"unknown operator: Boolean + Boolean",
			3,
		},
		{"foobar", "identifier not found: foobar", 1},
		{`"Hello" - "World"`, "unknown operator: String - String", 1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			fmt.Println(tt)
			t.Errorf("no error object returned. got=%T (%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}

		if errObj.Line() != tt.expectedLine {
			t.Errorf("wrong line. expected=%d, got=%d", tt.expectedLine, errObj.Line())
		}
	}
}

func TestVarDeclarationStmts(t *testing.T) {
	tests := []struct {
		source   string
		expected int64
	}{
		{"mut a = 5; a;", 5},
		{"const a = 5 * 5; a;", 25},
		{"mut a = 5; mut b = a; b;", 5},
		{"const a = 5; mut b = 5; const c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.source), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func(x) { x + 2; };"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"const identity = func(x) { x; }; identity(5);", 5},
		{"const identity = func(x) { return x; }; identity(5);", 5},
		{"const double = func(x) { x * 2; }; double(5);", 10},
		{"const add = func(x, y) { x + y; }; add(5, 5);", 10},
		{"const add = func(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"func(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
   const newAdder = func(x) {
     func(y) { x + y };
};
   const addTwo = newAdder(2);
   addTwo(2);`
	testIntegerObject(t, testEval(input), 4)
}

func TestVariableAssignment(t *testing.T) {
	tests := []struct {
		source   string
		expected int64
	}{
		{"mut y = 5; y = 6; y;", 6},
		{
			`
			mut x = 5;
			const fn = func () { x = 7; }
			fn();
			x;
			`,
			7,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.source)
		if !testIntegerObject(t, evaluated, tt.expected) {
			return
		}
	}
}

func TestStringLiteral(t *testing.T) {
	source := `"Hello, World!"`
	evaluated := testEval(source)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello, World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	source := `"Hello, " + "World!"`
	evaluated := testEval(source)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello, World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltInFunction(t *testing.T) {
	tests := []struct {
		source   string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` of wrong type. got=Integer"},
		{`len("one", "two")`, "`len` expects one argument"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.source)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}
	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testEval(source string) object.Object {
	l := lexer.New(source)
	p := parser.New(l)

	program := p.ParseProgram()
	scope := object.NewScope()
	return Eval(program, scope)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not *object.Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not *object.Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}
