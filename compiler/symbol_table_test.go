package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Scope: "GLOBAL", Index: 0, IsConstant: true},
		"b": {Name: "b", Scope: "GLOBAL", Index: 1, IsConstant: false},
	}

	globalScope := NewSymbolTable()

	a := globalScope.DefineImmutable("a")
	if a != expected["a"] {
		t.Errorf("a = %v, expected %v", a, expected["a"])
	}

	b := globalScope.DefineMutable("b")
	if b != expected["b"] {
		t.Errorf("b = %v, expected %v", b, expected["b"])
	}

}

func TestResolveGlobal(t *testing.T) {
	globalScope := NewSymbolTable()
	globalScope.DefineImmutable("a")
	globalScope.DefineMutable("b")

	expected := []Symbol{
		{Name: "a", Scope: "GLOBAL", Index: 0, IsConstant: true},
		{Name: "b", Scope: "GLOBAL", Index: 1, IsConstant: false},
	}

	for _, sym := range expected {
		result, ok := globalScope.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
		}

		if result != sym {
			t.Errorf("result = %v, expected %v", result, sym)
		}
	}
}
