package object

import (
	"QuonkScript/ast"
	"bytes"
	"fmt"
	"hash/fnv"

	"strings"
)

type ObjectType string
type BuiltInFunction func(line int, args ...Object) Object

const (
	IntegerObj     ObjectType = "Integer"
	BooleanObj     ObjectType = "Boolean"
	NullObj        ObjectType = "Null"
	ReturnValueObj ObjectType = "ReturnValue"
	ErrorObj       ObjectType = "Error"
	VariableObj    ObjectType = "Variable"
	FunctionObj    ObjectType = "Function"
	StringObj      ObjectType = "String"
	BuiltInObj     ObjectType = "BuiltIn"
	ArrayObj       ObjectType = "Array"
	HashObj        ObjectType = "Hash"
)

type (
	Object interface {
		Type() ObjectType
		Inspect() string
		Line() int
	}

	Hashable interface {
		HashKey() HashKey
	}
)

type (
	Integer struct {
		Value     int64
		TokenLine int
	}

	Boolean struct {
		Value     bool
		TokenLine int
	}

	Null struct {
		TokenLine int
	}

	ReturnValue struct {
		Value     Object
		TokenLine int
	}

	Error struct {
		Message    string
		OriginLine int
	}

	Variable struct {
		Value     Object
		Constant  bool
		TokenLine int
	}

	Function struct {
		Parameters []*ast.Identifier
		Body       *ast.BlockStmt
		Scope      *Scope
		TokenLine  int
	}

	String struct {
		Value     string
		TokenLine int
	}

	BuiltIn struct {
		Fn        BuiltInFunction
		TokenLine int
	}

	Array struct {
		Elements  []Object
		TokenLine int
	}

	HashKey struct {
		Type  ObjectType
		Value uint64
	}

	HashPair struct {
		Key   Object
		Value Object
	}

	Hash struct {
		Pairs     map[HashKey]HashPair
		TokenLine int
	}
)

func (i *Integer) Type() ObjectType {
	return IntegerObj
}

func (b *Boolean) Type() ObjectType {
	return BooleanObj
}

func (n *Null) Type() ObjectType {
	return NullObj
}

func (r *ReturnValue) Type() ObjectType {
	return ReturnValueObj
}

func (e *Error) Type() ObjectType {
	return ErrorObj
}

func (v *Variable) Type() ObjectType {
	return VariableObj
}

func (f *Function) Type() ObjectType {
	return FunctionObj
}

func (b *BuiltIn) Type() ObjectType {
	return BuiltInObj
}

func (s *String) Type() ObjectType {
	return StringObj
}

func (a *Array) Type() ObjectType {
	return ArrayObj
}

func (h *Hash) Type() ObjectType {
	return HashObj
}

func (i *Integer) Line() int {
	return i.TokenLine
}

func (b *Boolean) Line() int {
	return b.TokenLine
}

func (n *Null) Line() int {
	return n.TokenLine
}

func (r *ReturnValue) Line() int {
	return r.TokenLine
}

func (e *Error) Line() int {
	return e.OriginLine
}

func (f *Function) Line() int {
	return f.TokenLine
}

func (s *String) Line() int {
	return s.TokenLine
}

func (b *BuiltIn) Line() int {
	return b.TokenLine
}

func (a *Array) Line() int {
	return a.TokenLine
}

func (h *Hash) Line() int {
	return h.TokenLine
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

func (n *Null) Inspect() string {
	return "null"
}

func (e *Error) Inspect() string {
	return fmt.Sprintf("Honk! Error: %s on line %d", e.Message, e.Line())
}

func (v *Variable) Inspect() string {
	return v.Value.Inspect()
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("func(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

func (s *String) Inspect() string {
	return s.Value
}

func (b *BuiltIn) Inspect() string {
	return "builtin function"
}

func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// HashKey functions
func (b *Boolean) HashKey() HashKey {
	var val uint64

	if b.Value {
		val = 1
	} else {
		val = 0
	}

	return HashKey{Type: b.Type(), Value: val}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
