package types

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
