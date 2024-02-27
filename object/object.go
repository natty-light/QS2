package object

import "fmt"

type ObjectType string

const (
	IntegerObj ObjectType = "Integer"
	BooleanObj ObjectType = "Boolean"
	NullObj    ObjectType = "Null"
)

type RuntimeValue interface {
	Type() ObjectType
	Inspect() string
}

type (
	Integer struct {
		Value int64
	}

	Boolean struct {
		Value bool
	}

	Null struct{}
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

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (n *Null) Inspect() string {
	return "null"
}
