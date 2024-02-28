package object

import "fmt"

func NewScope() *Scope {
	s := make(map[string]Variable)

	return &Scope{store: s, outer: nil}
}

func NewEnclosedScope(outer *Scope) *Scope {
	s := NewScope()
	s.outer = outer

	return s
}

type Scope struct {
	store map[string]Variable
	outer *Scope
}

func (s *Scope) Get(name string) (Variable, bool, bool) {
	obj, ok := s.store[name]
	fromOuter := false
	if !ok && s.outer != nil {
		obj, ok, _ = s.outer.Get(name)
		fromOuter = true
	}
	return obj, ok, fromOuter
}

func (s *Scope) Set(name string, val Object, isConst bool) Object {

	existing, ok, fromOuter := s.Get(name)

	if !ok || !existing.Constant || fromOuter {
		// if the variable doesn't exist, or its not constant, or its from the parent scope
		s.store[name] = Variable{Value: val, Constant: isConst, TokenLine: val.Line()}
		return val
	} else {
		return newError(val.Line(), "attempt to assign a value to constant variable %s", name)
	}
}

func newError(line int, format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...), OriginLine: line}
}
