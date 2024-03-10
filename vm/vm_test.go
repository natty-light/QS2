package vm

import (
	"fmt"
	"quonk/ast"
	"quonk/compiler"
	"quonk/lexer"
	"quonk/object"
	"quonk/parser"
	"testing"
)

type vmTestCase struct {
	source   string
	expected interface{}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		// {"5 + 5 + 5 + 5 - 10", 10},
		// {"2 * 2 * 2 * 2 * 2", 32},
		// {"5 * 2 + 10", 20},
		// {"5 + 2 * 10", 25},
		// {"5 * (2 + 10)", 60},
	}

	runVmTests(t, tests)
}

func TestFLoatArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1.0", 1.0},
		{"2.0", 2.0},
		{"1.0 + 2.0", 3.0},
		{"1.0 * 2.0", 2.0},
		{"4.0 / 2.0", 2.0},
		{"50.0 / 2.0 * 2.0 + 10.0 - 5.0", 55.0},
		{"5.0 + 5.0 + 5.0 + 5.0 - 10.0", 10.0},
		{"2.0 * 2.0 * 2.0 * 2.0 * 2.0", 32.0},
		{"5.0 * 2.0 + 10.0", 20.0},
		{"5.0 + 2.0 * 10.0", 25.0},
		{"5.0 * (2.0 + 10.0)", 60.0},
	}

	runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
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
		{"1 < 2 == true", true},
		{"1 < 2 == false", false},
		{"1 > 2 == true", false},
		{"1 > 2 == false", true},
		{"1 < 2 == 2 > 1", true},
		{"1 < 2 == 2 < 1", false},
		{"1 > 2 == 2 > 1", false},
		{"1 > 2 == 2 < 1", true},
		{"(1 < 2) != true", false},
		{"(1 < 2) != false", true},
		{"(1 > 2) != true", true},
		{"(1 > 2) != false", false},
		{"(1 < 2) != (2 > 1)", false},
		{"(1 < 2) != (2 < 1)", true},
		{"(1 > 2) != (2 > 1)", true},
		{"(1 > 2) != (2 < 1)", false},
		{"true && true", true},
		{"true && false", false},
		{"false && true", false},
		{"false && false", false},
		{"true || true", true},
		{"true || false", true},
		{"false || true", true},
		{"false || false", false},
		{"1 > 2 && 1 < 2", false},
		{"1 > 2 || 1 < 2", true},
		{"1 < 2 && 1 < 2", true},
		{"true == true && true == true", true},
		{"true == true && true == false", false},
		{"true == true || true == false", true},
		{"true == false || true == false", false},
		{"true != true && true != true", false},
		{"true != true && true != false", false},
		{"true != true || true != false", true},
		{"true != false || true != false", true},
	}

	runVmTests(t, tests)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.source)

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

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case float64:
		err := testFloatObject(expected, actual)
		if err != nil {
			t.Errorf("testFloatObject failed: %s", err)
		}
	case float32:
		err := testFloatObject(float64(expected), actual)
		if err != nil {
			t.Errorf("testFloatObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(expected, actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	}
}

func parse(source string) *ast.Program {
	l := lexer.New(source)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}

	return nil
}

func testFloatObject(expected float64, actual object.Object) error {
	result, ok := actual.(*object.Float)
	if !ok {
		return fmt.Errorf("object is not Float. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
	}

	return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}

	return nil
}
