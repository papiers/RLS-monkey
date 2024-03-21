package object

// NewEnvironment 创建环境对象
func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
		outer: nil,
	}
}

// NewEnclosedEnvironment 创建封闭的环境对象
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Environment 存储变量名和变量的映射关系
type Environment struct {
	store map[string]Object
	outer *Environment
}

// Get 获取变量
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set 设置变量
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return e.store[name]
}
