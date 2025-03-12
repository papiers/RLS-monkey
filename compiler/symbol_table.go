package compiler

// SymbolScope 符号作用域
type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
	LocalScope  SymbolScope = "LOCAL"
)

// Symbol 符号
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable 符号表
type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int
}

// NewSymbolTable 创建符号表
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store: make(map[string]Symbol),
	}
}

// Define 定义符号
func (st *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{
		Name:  name,
		Scope: GlobalScope,
		Index: st.numDefinitions,
	}
	st.store[name] = symbol
	st.numDefinitions++
	return symbol
}

// Resolve 解析符号
func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := st.store[name]
	return symbol, ok
}
