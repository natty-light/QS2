package compiler

type SymbolScopes string

const (
	GlobalScope SymbolScopes = "GLOBAL"
	LocalScope  SymbolScopes = "LOCAL"
)

type Symbol struct {
	Name       string
	Scope      SymbolScopes
	Index      int
	IsConstant bool
}

type SymbolTable struct {
	Outer *SymbolTable

	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func (s *SymbolTable) DefineMutable(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions, IsConstant: false}
	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

func (s *SymbolTable) DefineImmutable(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions, IsConstant: true}
	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := s.store[name]
	if !ok && s.Outer != nil {
		symbol, ok = s.Outer.Resolve(name)
	}
	return symbol, ok
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}
