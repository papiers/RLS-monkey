package object

import "fmt"

// Builtins 保存内置函数
var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		"len",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1",
						len(args))
				}
				switch arg := args[0].(type) {
				case *Array:
					return &Integer{Value: int64(len(arg.Elements))}
				case *String:
					return &Integer{Value: int64(len(arg.Value))}
				default:
					return newError("argument to `len` not supported, got %s",
						args[0].Type())
				}
			},
		},
	},
	{
		"puts",
		&Builtin{
			Fn: func(args ...Object) Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return nil
			},
		},
	},
	{
		"first",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return &Error{
						Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)),
					}
				}
				switch arg := args[0].(type) {
				case *Array:
					if len(arg.Elements) > 0 {
						return arg.Elements[0]
					}
					return nil
				default:
					return &Error{
						Message: fmt.Sprintf("argument to `first` must be Array, got %s", arg.Type()),
					}
				}
			}},
	},
	{
		"last",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return &Error{
						Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)),
					}
				}
				switch arg := args[0].(type) {
				case *Array:
					l := len(arg.Elements)
					if l > 0 {
						return arg.Elements[l-1]
					}
					return nil
				default:
					return &Error{
						Message: fmt.Sprintf("argument to `last` must be Array, got %s", arg.Type()),
					}
				}
			},
		},
	},
	{
		"rest",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return &Error{
						Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)),
					}
				}
				switch arg := args[0].(type) {
				case *Array:
					l := len(arg.Elements)
					if l > 0 {
						newElements := make([]Object, l-1)
						copy(newElements, arg.Elements[1:])
						return &Array{Elements: newElements}
					}
					return nil
				default:
					return &Error{
						Message: fmt.Sprintf("argument to `rest` must be Array, got %s", arg.Type()),
					}
				}
			},
		},
	},
	{
		"push",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return &Error{
						Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args)),
					}
				}
				switch arg := args[0].(type) {
				case *Array:
					l := len(arg.Elements)
					newElements := make([]Object, l+1, l+1)
					copy(newElements, arg.Elements)
					newElements[l] = args[1]
					return &Array{Elements: newElements}
				default:
					return &Error{
						Message: fmt.Sprintf("argument to `push` must be Array, got %s", arg.Type()),
					}
				}
			},
		},
	},
	{
		"",
		&Builtin{},
	},
}

// newError 返回一个错误对象
func newError(format string, a ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

// GetBuiltinByName 根据名字获取内置函数
func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}
