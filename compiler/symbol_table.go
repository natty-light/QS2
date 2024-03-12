package compiler

type SymbolScopes string

const (
	GlobalScope SymbolScopes = "GLOBAL"
)

type Symbol struct {
	Name       string
	Scope      SymbolScopes
	Index      int
	IsConstant bool
}

type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func (s *SymbolTable) DefineMutable(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions, Scope: GlobalScope, IsConstant: false}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

func (s *SymbolTable) DefineImmutable(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions, Scope: GlobalScope, IsConstant: true}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := s.store[name]
	return symbol, ok
}
