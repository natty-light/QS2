package object

import "fmt"

type ObjectType string

const (
	IntegerObj     ObjectType = "Integer"
	BooleanObj     ObjectType = "Boolean"
	NullObj        ObjectType = "Null"
	ReturnValueObj ObjectType = "ReturnValue"
)

type Object interface {
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

	ReturnValue struct {
		Value Object
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
