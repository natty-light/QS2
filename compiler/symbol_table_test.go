package compiler

import (
	"quonk/types"
	"testing"
)

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Scope: "GLOBAL", Index: 0, IsConstant: true, Type: &types.Int{}},
		"b": {Name: "b", Scope: "GLOBAL", Index: 1, IsConstant: false, Type: &types.Int{}},
		"c": {Name: "c", Scope: "LOCAL", Index: 0, IsConstant: true, Type: &types.Int{}},
		"d": {Name: "d", Scope: "LOCAL", Index: 1, IsConstant: false, Type: &types.Int{}},
		"e": {Name: "e", Scope: "LOCAL", Index: 0, IsConstant: true, Type: &types.Int{}},
		"f": {Name: "f", Scope: "LOCAL", Index: 1, IsConstant: false, Type: &types.Int{}},
	}

	globalScope := NewSymbolTable()

	a := globalScope.DefineImmutable("a", &types.Int{})
	if a != expected["a"] {
		t.Errorf("a = %v, expected %v", a, expected["a"])
	}

	b := globalScope.DefineMutable("b", &types.Int{})
	if b != expected["b"] {
		t.Errorf("b = %v, expected %v", b, expected["b"])
	}

	firstLocal := NewEnclosedSymbolTable(globalScope)
	c := firstLocal.DefineImmutable("c", &types.Int{})
	if c != expected["c"] {
		t.Errorf("c = %v, expected %v", c, expected["c"])
	}

	d := firstLocal.DefineMutable("d", &types.Int{})
	if d != expected["d"] {
		t.Errorf("d = %v, expected %v", d, expected["d"])
	}

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	e := secondLocal.DefineImmutable("e", &types.Int{})
	if e != expected["e"] {
		t.Errorf("e = %v, expected %v", e, expected["e"])
	}

	f := secondLocal.DefineMutable("f", &types.Int{})
	if f != expected["f"] {
		t.Errorf("f = %v, expected %v", f, expected["f"])
	}

}

func TestResolveGlobal(t *testing.T) {
	globalScope := NewSymbolTable()
	globalScope.DefineImmutable("a", &types.Int{})
	globalScope.DefineMutable("b", &types.Int{})

	expected := []Symbol{
		{Name: "a", Scope: "GLOBAL", Index: 0, IsConstant: true, Type: &types.Int{}},
		{Name: "b", Scope: "GLOBAL", Index: 1, IsConstant: false, Type: &types.Int{}},
	}

	for _, sym := range expected {
		result, _, ok := globalScope.Resolve(sym.Name)
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
	globalScope.DefineImmutable("a", &types.Int{})
	globalScope.DefineMutable("b", &types.Int{})

	firstLocal := NewEnclosedSymbolTable(globalScope)
	firstLocal.DefineImmutable("c", &types.Int{})
	firstLocal.DefineMutable("d", &types.Int{})

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.DefineImmutable("e", &types.Int{})
	secondLocal.DefineMutable("f", &types.Int{})

	tests := []struct {
		table           *SymbolTable
		expectedSymbols []Symbol
	}{
		{
			firstLocal,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0, IsConstant: true, Type: &types.Int{}},
				{Name: "b", Scope: GlobalScope, Index: 1, IsConstant: false, Type: &types.Int{}},
				{Name: "c", Scope: LocalScope, Index: 0, IsConstant: true, Type: &types.Int{}},
				{Name: "d", Scope: LocalScope, Index: 1, IsConstant: false, Type: &types.Int{}},
			},
		},
		{
			secondLocal,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0, IsConstant: true, Type: &types.Int{}},
				{Name: "b", Scope: GlobalScope, Index: 1, IsConstant: false, Type: &types.Int{}},
				{Name: "e", Scope: LocalScope, Index: 0, IsConstant: true, Type: &types.Int{}},
				{Name: "f", Scope: LocalScope, Index: 1, IsConstant: false, Type: &types.Int{}},
			},
		},
	}

	for _, tt := range tests {
		for _, sym := range tt.expectedSymbols {
			result, _, ok := tt.table.Resolve(sym.Name)
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
	globalScope.DefineImmutable("a", &types.Int{})
	globalScope.DefineMutable("b", &types.Int{})

	local := NewEnclosedSymbolTable(globalScope)
	local.DefineImmutable("c", &types.Int{})
	local.DefineMutable("d", &types.Int{})

	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0, IsConstant: true, Type: &types.Int{}},
		{Name: "b", Scope: GlobalScope, Index: 1, IsConstant: false, Type: &types.Int{}},
		{Name: "c", Scope: LocalScope, Index: 0, IsConstant: true, Type: &types.Int{}},
		{Name: "d", Scope: LocalScope, Index: 1, IsConstant: false, Type: &types.Int{}},
	}

	for _, sym := range expected {
		result, _, ok := local.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
		}
	}
}

func TestDefineResolveBuiltins(t *testing.T) {
	global := NewSymbolTable()
	firstLocal := NewEnclosedSymbolTable(global)
	secondLocal := NewEnclosedSymbolTable(firstLocal)

	expected := []Symbol{
		{Name: "a", Scope: BuiltinScope, Index: 0, IsConstant: true, Type: &types.Int{}},
		{Name: "b", Scope: BuiltinScope, Index: 1, IsConstant: true, Type: &types.Int{}},
		{Name: "c", Scope: BuiltinScope, Index: 2, IsConstant: true, Type: &types.Int{}},
		{Name: "d", Scope: BuiltinScope, Index: 3, IsConstant: true, Type: &types.Int{}},
	}

	for i, sym := range expected {
		global.DefineBuiltin(i, sym.Name)
	}

	for _, table := range []*SymbolTable{global, firstLocal, secondLocal} {
		for _, sym := range expected {
			result, _, ok := table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
			}

			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
			}
		}
	}
}

func TestResolveFree(t *testing.T) {
	global := NewSymbolTable()
	global.DefineImmutable("a", &types.Int{})
	global.DefineImmutable("b", &types.Int{})

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.DefineImmutable("c", &types.Int{})
	firstLocal.DefineImmutable("d", &types.Int{})

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.DefineImmutable("e", &types.Int{})
	secondLocal.DefineImmutable("f", &types.Int{})

	tests := []struct {
		table               *SymbolTable
		expectedSymbols     []Symbol
		expectedFreeSymbols []Symbol
	}{
		{
			firstLocal,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0, IsConstant: true, Type: &types.Int{}},
				{Name: "b", Scope: GlobalScope, Index: 1, IsConstant: true, Type: &types.Int{}},
				{Name: "c", Scope: LocalScope, Index: 0, IsConstant: true, Type: &types.Int{}},
				{Name: "d", Scope: LocalScope, Index: 1, IsConstant: true, Type: &types.Int{}},
			},
			[]Symbol{},
		},
		{
			secondLocal,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0, IsConstant: true, Type: &types.Int{}},
				{Name: "b", Scope: GlobalScope, Index: 1, IsConstant: true, Type: &types.Int{}},
				{Name: "c", Scope: FreeScope, Index: 0, IsConstant: true, Type: &types.Int{}},
				{Name: "d", Scope: FreeScope, Index: 1, IsConstant: true, Type: &types.Int{}},
				{Name: "e", Scope: LocalScope, Index: 0, IsConstant: true, Type: &types.Int{}},
				{Name: "f", Scope: LocalScope, Index: 1, IsConstant: true, Type: &types.Int{}},
			},
			[]Symbol{
				{Name: "c", Scope: LocalScope, Index: 0, IsConstant: true, Type: &types.Int{}},
				{Name: "d", Scope: LocalScope, Index: 1, IsConstant: true, Type: &types.Int{}},
			},
		},
	}

	for _, tt := range tests {
		for _, sym := range tt.expectedSymbols {
			result, _, ok := tt.table.Resolve(sym.Name)

			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}

			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
			}
		}

		if len(tt.table.FreeSymbols) != len(tt.expectedFreeSymbols) {
			t.Errorf("wrong number of free symbols. got=%d, want=%d", len(tt.table.FreeSymbols), len(tt.expectedFreeSymbols))
			continue
		}

		for i, sym := range tt.expectedFreeSymbols {
			result := tt.table.FreeSymbols[i]
			if result != sym {
				t.Errorf("wrong free symbol. got=%+v, want=%+v", result, sym)
			}
		}
	}
}

func TestResolveUnresolvableFree(t *testing.T) {
	global := NewSymbolTable()
	global.DefineImmutable("a", &types.Int{})

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.DefineImmutable("c", &types.Int{})

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.DefineImmutable("e", &types.Int{})
	secondLocal.DefineImmutable("f", &types.Int{})

	expected := []Symbol{
		{"a", GlobalScope, 0, true, &types.Int{}},
		{"c", FreeScope, 0, true, &types.Int{}},
		{"e", LocalScope, 0, true, &types.Int{}},
		{"f", LocalScope, 1, true, &types.Int{}},
	}

	for _, sym := range expected {
		result, _, ok := secondLocal.Resolve(sym.Name)

		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
		}
	}

	expectedUnresolvable := []string{
		"b",
		"d",
	}

	for _, name := range expectedUnresolvable {
		_, _, ok := secondLocal.Resolve(name)
		if ok {
			t.Errorf("name %s resolved, but was expected not to", name)
		}
	}
}

func TestDefineAndResolveFunctionName(t *testing.T) {
	g := NewSymbolTable()
	g.DefineFunctionName("a")

	expected := Symbol{Name: "a", Scope: FunctionScope, Index: 0, IsConstant: false}

	result, _, ok := g.Resolve(expected.Name)
	if !ok {
		t.Fatalf("function name %s not resolvable", expected.Name)
	}

	if result != expected {
		t.Errorf("expected %s to resolve to %+v, got=%+v", expected.Name, expected, result)
	}
}

func TestShadowingFunctionName(t *testing.T) {
	g := NewSymbolTable()
	g.DefineFunctionName("a")
	g.DefineImmutable("a", &types.Int{})

	expected := Symbol{Name: "a", Scope: GlobalScope, Index: 0, IsConstant: true}

	result, _, ok := g.Resolve(expected.Name)
	if !ok {
		t.Fatalf("function name %s not resolvable", expected.Name)
	}

	if result != expected {
		t.Errorf("expected %s to resolve to %+v, got=%+v", expected.Name, expected, result)
	}
}
