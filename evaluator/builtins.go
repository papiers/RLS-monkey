package evaluator

import (
	"fmt"

	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{
					Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)),
				}
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return &object.Error{
					Message: fmt.Sprintf("argument to `len` not supported. got=%s", arg.Type()),
				}
			}
		},
	},
	"put": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{
					Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args)),
				}
			}
			switch arg := args[0].(type) {
			case *object.Array:
				arg.Elements = append(arg.Elements, args[1])
			default:
				return &object.Error{
					Message: fmt.Sprintf("argument to `put` must be Array, got %s", arg.Type()),
				}
			}
			return &object.Null{}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{
					Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)),
				}
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[0]
				}
				return &object.Null{}
			default:
				return &object.Error{
					Message: fmt.Sprintf("argument to `first` must be Array, got %s", arg.Type()),
				}
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{
					Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)),
				}
			}
			switch arg := args[0].(type) {
			case *object.Array:
				l := len(arg.Elements)
				if l > 0 {
					return arg.Elements[l-1]
				}
				return &object.Null{}
			default:
				return &object.Error{
					Message: fmt.Sprintf("argument to `last` must be Array, got %s", arg.Type()),
				}
			}
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{
					Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)),
				}
			}
			switch arg := args[0].(type) {
			case *object.Array:
				l := len(arg.Elements)
				if l > 0 {
					newElements := make([]object.Object, l-1)
					copy(newElements, arg.Elements[1:])
					return &object.Array{Elements: newElements}
				}
				return &object.Null{}
			default:
				return &object.Error{
					Message: fmt.Sprintf("argument to `rest` must be Array, got %s", arg.Type()),
				}
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{
					Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args)),
				}
			}
			switch arg := args[0].(type) {
			case *object.Array:
				arg.Elements = append(arg.Elements, args[1])
				return arg
			default:
				return &object.Error{
					Message: fmt.Sprintf("argument to `push` must be Array, got %s", arg.Type()),
				}
			}
		},
	},
}
