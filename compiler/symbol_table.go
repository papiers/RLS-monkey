package compiler

// SymbolScope 符号作用域
type SymbolScope string

const (
	GlobalScope   SymbolScope = "GLOBAL"
	LocalScope    SymbolScope = "LOCAL"
	BuiltinScope  SymbolScope = "BUILTIN"
	FreeScope     SymbolScope = "FREE"
	FunctionScope SymbolScope = "FUNCTION"
)

// Symbol 符号
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable 符号表
type SymbolTable struct {
	Outer *SymbolTable

	store          map[string]Symbol
	numDefinitions int

	FreeSymbols []Symbol
}

// NewSymbolTable 创建符号表
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store:       make(map[string]Symbol),
		FreeSymbols: []Symbol{},
	}
}

// NewEnclosedSymbolTable 创建封闭的符号表
func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

// Define 定义符号
func (st *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{
		Name:  name,
		Index: st.numDefinitions,
	}
	if st.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}
	st.store[name] = symbol
	st.numDefinitions++
	return symbol
}

// Resolve 解析符号
func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := st.store[name]
	if !ok && st.Outer != nil {
		symbol, ok = st.Outer.Resolve(name)
		if ok && symbol.Scope != GlobalScope && symbol.Scope != BuiltinScope {
			symbol = st.DefineFree(symbol)
			ok = true
		}
	}
	return symbol, ok
}

// DefineBuiltin 定义内置符号
func (st *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{
		Name:  name,
		Index: index,
		Scope: BuiltinScope,
	}
	st.store[name] = symbol
	return symbol
}

// DefineFree 定义自由变量
func (st *SymbolTable) DefineFree(original Symbol) Symbol {
	st.FreeSymbols = append(st.FreeSymbols, original)
	symbol := Symbol{
		Name:  original.Name,
		Index: len(st.FreeSymbols) - 1,
		Scope: FreeScope,
	}
	st.store[original.Name] = symbol
	return symbol
}

// DefineFunctionName 定义函数名
func (st *SymbolTable) DefineFunctionName(name string) Symbol {
	symbol := Symbol{
		Name:  name,
		Scope: FunctionScope,
	}
	st.store[name] = symbol
	return symbol
}
