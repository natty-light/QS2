package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Scope: "GLOBAL", Index: 0, IsConstant: true},
		"b": {Name: "b", Scope: "GLOBAL", Index: 1, IsConstant: false},
		"c": {Name: "c", Scope: "LOCAL", Index: 0, IsConstant: true},
		"d": {Name: "d", Scope: "LOCAL", Index: 1, IsConstant: false},
		"e": {Name: "e", Scope: "LOCAL", Index: 0, IsConstant: true},
		"f": {Name: "f", Scope: "LOCAL", Index: 1, IsConstant: false},
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

	firstLocal := NewEnclosedSymbolTable(globalScope)
	c := firstLocal.DefineImmutable("c")
	if c != expected["c"] {
		t.Errorf("c = %v, expected %v", c, expected["c"])
	}

	d := firstLocal.DefineMutable("d")
	if d != expected["d"] {
		t.Errorf("d = %v, expected %v", d, expected["d"])
	}

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	e := secondLocal.DefineImmutable("e")
	if e != expected["e"] {
		t.Errorf("e = %v, expected %v", e, expected["e"])
	}

	f := secondLocal.DefineMutable("f")
	if f != expected["f"] {
		t.Errorf("f = %v, expected %v", f, expected["f"])
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
			t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
		}
	}
}

func TestResolveNestedLocal(t *testing.T) {
	globalScope := NewSymbolTable()
	globalScope.DefineImmutable("a")
	globalScope.DefineMutable("b")

	firstLocal := NewEnclosedSymbolTable(globalScope)
	firstLocal.DefineImmutable("c")
	firstLocal.DefineMutable("d")

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.DefineImmutable("e")
	secondLocal.DefineMutable("f")

	tests := []struct {
		table           *SymbolTable
		expectedSymbols []Symbol
	}{
		{
			firstLocal,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0, IsConstant: true},
				{Name: "b", Scope: GlobalScope, Index: 1, IsConstant: false},
				{Name: "c", Scope: LocalScope, Index: 0, IsConstant: true},
				{Name: "d", Scope: LocalScope, Index: 1, IsConstant: false},
			},
		},
		{
			secondLocal,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0, IsConstant: true},
				{Name: "b", Scope: GlobalScope, Index: 1, IsConstant: false},
				{Name: "e", Scope: LocalScope, Index: 0, IsConstant: true},
				{Name: "f", Scope: LocalScope, Index: 1, IsConstant: false},
			},
		},
	}

	for _, tt := range tests {
		for _, sym := range tt.expectedSymbols {
			result, ok := tt.table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
			}

			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
			}
		}
	}
}

func TestResolveLocal(t *testing.T) {
	globalScope := NewSymbolTable()
	globalScope.DefineImmutable("a")
	globalScope.DefineMutable("b")

	local := NewEnclosedSymbolTable(globalScope)
	local.DefineImmutable("c")
	local.DefineMutable("d")

	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0, IsConstant: true},
		{Name: "b", Scope: GlobalScope, Index: 1, IsConstant: false},
		{Name: "c", Scope: LocalScope, Index: 0, IsConstant: true},
		{Name: "d", Scope: LocalScope, Index: 1, IsConstant: false},
	}

	for _, sym := range expected {
		result, ok := local.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
		}
	}
}
