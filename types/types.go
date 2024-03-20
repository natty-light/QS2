package types

import "bytes"

type DataType string

const (
	IntType   DataType = "IntType"
	BoolType  DataType = "BoolType"
	FloatType DataType = "FloatType"
	StrType   DataType = "StrType"
	ArrayType DataType = "ArrayType"
	HashType  DataType = "HashType"
	FuncType  DataType = "FuncType"
	NullType  DataType = "NullType"
)

type Type interface {
	Type() DataType
	String() string
}

type (
	Int   struct{}
	Bool  struct{}
	Float struct{}
	Str   struct{}
	Array struct {
		Element Type
	}
	Hash struct {
		Key   Type
		Value Type
	}

	Func struct {
		Parameters []Type
		Return     Type
	}
)

func (i *Int) Type() DataType {
	return IntType
}

func (b *Bool) Type() DataType {
	return BoolType
}

func (f *Float) Type() DataType {
	return FloatType
}

func (s *Str) Type() DataType {
	return StrType
}

func (a *Array) Type() DataType {
	return ArrayType
}

func (h *Hash) Type() DataType {
	return HashType
}

func (f *Func) Type() DataType {
	return FuncType
}

func (a *Array) ElementType() Type {
	return a.Element
}

func (h *Hash) KeyType() Type {
	return h.Key
}

func (h *Hash) ValueType() Type {
	return h.Value
}

func (f *Func) ReturnType() Type {
	return f.Return
}

func (f *Func) ParameterTypes() []Type {
	return f.Parameters
}

func (i *Int) String() string {
	return "int"
}

func (b *Bool) String() string {
	return "bool"
}

func (f *Float) String() string {
	return "float"
}

func (s *Str) String() string {
	return "str"
}

func (a *Array) String() string {
	return "[]" + a.Element.String()
}

func (h *Hash) String() string {
	return "{" + h.Key.String() + " " + h.Value.String() + "}"
}

func (f *Func) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	for i, p := range f.Parameters {
		out.WriteString(p.String())
		if i != len(f.Parameters)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(") -> ")
	out.WriteString(f.Return.String())

	return out.String()
}
