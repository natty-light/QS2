package evaluator

import (
	"bytes"
	"fmt"
	"quonk/object"
	"strings"
)

var builtIns = map[string]*object.BuiltIn{
	"len": {
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
	"first": {
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
	"last": {
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
	"rest": {
		Fn: func(line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(line, "`rest` expects one argument")
			}
			if args[0].Type() != object.ArrayObj {
				return newError(line, "argument of `rest` must be array type")
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElems := make([]object.Object, length-1)
				copy(newElems, arr.Elements[1:length])
				return &object.Array{Elements: newElems}
			}
			return NULL
		},
	},
	"append": {
		Fn: func(line int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(line, "`append` expects two arguments")
			}

			if args[0].Type() != object.ArrayObj {
				return newError(line, "first argument to `append` must be array")
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElems := make([]object.Object, length+1)
			copy(newElems, arr.Elements)
			newElems[length] = args[1]

			return &object.Array{Elements: newElems}
		},
	},
	"slice": {
		Fn: func(line int, args ...object.Object) object.Object {
			if len(args) != 3 {
				return newError(line, "`slice` expects three arguments")
			}

			if args[0].Type() != object.ArrayObj {
				return newError(line, "first argument to `slice` must be array")
			}
			if args[1].Type() != object.IntegerObj {
				return newError(line, "`start` argument to `slice` must be int")
			}

			if args[2].Type() != object.IntegerObj {
				return newError(line, "`end` argument to `slice` must be int")
			}

			arr := args[0].(*object.Array)
			start := args[1].(*object.Integer).Value
			end := args[2].(*object.Integer).Value
			arrLength := int64(len(arr.Elements) - 1)

			if start < 0 {
				start = 0
			}
			if end > arrLength {
				end = arrLength - 1
			}
			slicedLength := int64(end - start)

			newElems := make([]object.Object, slicedLength)
			copy(newElems, arr.Elements[start:end])

			return &object.Array{Elements: newElems}
		},
	},
	"print": {
		Fn: func(line int, args ...object.Object) object.Object {
			var out bytes.Buffer

			elems := make([]string, 0)
			for _, a := range args {
				elems = append(elems, a.Inspect())
			}
			out.WriteString(strings.Join(elems, " "))

			fmt.Println(out.String())
			return nil
		},
	},
	"keys": {
		Fn: func(line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(line, "`keys` expects one argument")
			}

			hash, ok := args[0].(*object.Hash)
			if !ok {
				return newError(line, "unknown argument type for `keys`: %T", args[0])
			}

			keys := make([]object.Object, 0)
			for key := range hash.Pairs {
				switch val := key.ObjectValue.(type) {
				case bool:
					keys = append(keys, &object.Boolean{Value: val})
				case string:
					keys = append(keys, &object.String{Value: val})
				case int64:
					keys = append(keys, &object.Integer{Value: val})
				}
			}

			return &object.Array{Elements: keys}
		},
	},
}
