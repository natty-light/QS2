package object

import "fmt"

func NewScope() *Scope {
	s := make(map[string]Variable)

	return &Scope{store: s}
}

type Scope struct {
	store map[string]Variable
}

func (s *Scope) Get(name string) (Variable, bool) {
	obj, ok := s.store[name]
	return obj, ok
}

func (s *Scope) Set(name string, val Object, isConst bool) Object {

	existing, ok := s.Get(name)

	if !ok && !existing.Constant {
		s.store[name] = Variable{Value: val, Constant: isConst, TokenLine: val.Line()}
		return val
	} else {
		return newError(val.Line(), "attempt to assign a value to constant variable %s", name)
	}

}

func newError(line int, format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...), OriginLine: line}
}
