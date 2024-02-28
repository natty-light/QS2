package evaluator

import (
	"QuonkScript/lexer"
	"QuonkScript/object"
	"QuonkScript/parser"
	"testing"
)

func TestEvalIntegerExpr(t *testing.T) {
	tests := []struct {
		source   string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.source)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
func testEval(source string) object.Object {
	l := lexer.New(source)
	p := parser.New(l)

	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not *object.Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", expected, result.Value)
		return false
	}

	return true
}
