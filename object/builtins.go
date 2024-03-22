package object

import (
	"bytes"
	"fmt"
	"strings"
)

var Builtins = []struct {
	Name    string
	BuiltIn *BuiltIn
}{
	{
		"len",
		&BuiltIn{
			Fn: func(line int, args ...Object) Object {
				if len(args) != 1 {
					return newError(line, "`len` expects one argument")
				}

				switch arg := args[0].(type) {
				case *String:
					return &Integer{Value: int64(len(arg.Value))}
				case *Array:
					return &Integer{Value: int64(len(arg.Elements))}
				default:
					return newError(line, "argument to `len` of wrong type. got=%s", args[0].Type())
				}
			},
		},
	},
	{
		"print",
		&BuiltIn{
			Fn: func(line int, args ...Object) Object {
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
	},
	{
		"first",
		&BuiltIn{
			Fn: func(line int, args ...Object) Object {
				if len(args) != 1 {
					return newError(line, "`first` expects a single argument. got=%d", len(args))
				}
				if args[0].Type() != ArrayObj {
					return newError(line, "argument to `first` must be an Array. got=%s", args[0].Type())
				}

				arr := args[0].(*Array)
				if len(arr.Elements) > 0 {
					return arr.Elements[0]
				}
				return nil
			},
		},
	},
	{
		"last",
		&BuiltIn{
			Fn: func(line int, args ...Object) Object {
				if len(args) != 1 {
					return newError(line, "`last` expects a single argument. got=%d", len(args))
				}
				if args[0].Type() != ArrayObj {
					return newError(line, "argument to `last` must be an Array. got=%s", args[0].Type())
				}

				arr := args[0].(*Array)
				length := len(arr.Elements)
				if length > 0 {
					return arr.Elements[length-1]
				}
				return nil
			},
		},
	},
	{
		"rest",
		&BuiltIn{
			Fn: func(line int, args ...Object) Object {
				if len(args) != 1 {
					return newError(line, "`rest` expects one argument")
				}
				if args[0].Type() != ArrayObj {
					return newError(line, "argument of `rest` must be array type")
				}
				arr := args[0].(*Array)
				length := len(arr.Elements)
				if length > 0 {
					newElems := make([]Object, length-1)
					copy(newElems, arr.Elements[1:length])
					return &Array{Elements: newElems}
				}
				return nil
			},
		},
	},
	{
		"append",
		&BuiltIn{
			Fn: func(line int, args ...Object) Object {
				if len(args) != 2 {
					return newError(line, "`append` expects two arguments")
				}

				if args[0].Type() != ArrayObj {
					return newError(line, "first argument to `append` must be array")
				}

				arr := args[0].(*Array)
				length := len(arr.Elements)

				newElems := make([]Object, length+1)
				copy(newElems, arr.Elements)
				newElems[length] = args[1]

				return &Array{Elements: newElems}
			},
		},
	},
	{
		"slice",
		&BuiltIn{
			Fn: func(line int, args ...Object) Object {
				if len(args) != 3 {
					return newError(line, "`slice` expects three arguments")
				}

				if args[0].Type() != ArrayObj {
					return newError(line, "first argument to `slice` must be array")
				}
				if args[1].Type() != IntegerObj {
					return newError(line, "`start` argument to `slice` must be int")
				}

				if args[2].Type() != IntegerObj {
					return newError(line, "`end` argument to `slice` must be int")
				}

				arr := args[0].(*Array)
				start := args[1].(*Integer).Value
				end := args[2].(*Integer).Value
				arrLength := int64(len(arr.Elements) - 1)

				if start < 0 {
					start = 0
				}
				if end > arrLength {
					end = arrLength - 1
				}
				slicedLength := int64(end - start)

				newElems := make([]Object, slicedLength)
				copy(newElems, arr.Elements[start:end])

				return &Array{Elements: newElems}
			},
		},
	},
	{
		"keys",
		&BuiltIn{
			Fn: func(line int, args ...Object) Object {
				if len(args) != 1 {
					return newError(line, "`keys` expects one argument")
				}

				hash, ok := args[0].(*Hash)
				if !ok {
					return newError(line, "unknown argument type for `keys`: %T", args[0])
				}

				keys := make([]Object, 0)
				for key := range hash.Pairs {
					switch val := key.ObjectValue.(type) {
					case bool:
						keys = append(keys, &Boolean{Value: val})
					case string:
						keys = append(keys, &String{Value: val})
					case int64:
						keys = append(keys, &Integer{Value: val})
					}
				}

				return &Array{Elements: keys}
			},
		},
	},
	{
		"values",
		&BuiltIn{
			Fn: func(line int, args ...Object) Object {
				if len(args) != 1 {
					return newError(line, "`values` expects one argument")
				}

				hash, ok := args[0].(*Hash)
				if !ok {
					return newError(line, "unknown argument type for `values`: %T", args[0])
				}

				values := make([]Object, 0)
				for _, pair := range hash.Pairs {
					values = append(values, pair.Value)
				}

				return &Array{Elements: values}
			},
		},
	},
}

func GetBuiltInByName(name string) *BuiltIn {
	for _, bi := range Builtins {
		if bi.Name == name {
			return bi.BuiltIn
		}
	}

	return nil
}
