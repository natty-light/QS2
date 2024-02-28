package object

import "fmt"

type ObjectType string

const (
	IntegerObj     ObjectType = "Integer"
	BooleanObj     ObjectType = "Boolean"
	NullObj        ObjectType = "Null"
	ReturnValueObj ObjectType = "ReturnValue"
	ErrorObj       ObjectType = "Error"
	VariableObj    ObjectType = "Variable"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	Line() int
}

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
