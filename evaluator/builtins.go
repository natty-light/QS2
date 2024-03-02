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
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError(line, "argument to `len` of wrong type. got=%s", arg.Type())
			}
		},
	},
	"first": &object.BuiltIn{
		Fn: func(line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(line, "`first` expects a single argument. got=%d", len(args))
			}
			if args[0].Type() != object.ArrayObj {
				return newError(line, "argument to `first` must be an Array. got=%s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},
	"last": &object.BuiltIn{
		Fn: func(line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(line, "`last` expects a single argument. got=%d", len(args))
			}
			if args[0].Type() != object.ArrayObj {
				return newError(line, "argument to `last` must be an Array. got=%s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},
	// This has a dependecy cycle caused by calling Eval inside of applyFunction
	//"map": &object.BuiltIn{
	//	Fn: func(line int, args ...object.Object) object.Object {
	//		if len(args) != 2 {
	//			return newError(line, "`map` expects 2 arguments. got=%d", len(args))
	//		}
	//
	//		if args[0].Type() != object.ArrayObj {
	//			return newError(line, "`map` expects array as first argument")
	//		}
	//		arr := args[0].(*object.Array)
	//
	//		if args[1].Type() != object.FunctionObj {
	//			return newError(line, "`map` expects callback as first argument")
	//		}
	//		// Callback will have its own scope
	//		callback := args[1].(*object.Function)
	//		ret := &object.Array{TokenLine: line, Elements: make([]object.Object, 0)}
	//
	//		for _, e := range arr.Elements {
	//			res := applyFunction(callback, []object.Object{e}, line)
	//			ret.Elements = append(ret.Elements, res)
	//		}
	//		return ret
	//	},
	//},
}
