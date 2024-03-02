package evaluator

import "QuonkScript/object"

var builtIns = map[string]*object.BuiltIn{
	"len": &object.BuiltIn{
		Fn: func(line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(line, "`len` expects one argument")
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError(line, "argument to `len` of wrong type. got=%s", arg.Type())
			}
		},
	},
}
