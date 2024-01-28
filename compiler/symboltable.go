package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
	LocalScope  SymbolScope = "LOCAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Outer *SymbolTable // 用于访问外层符号表

	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions}

	if s.Outer != nil {
		symbol.Scope = LocalScope
	} else {
		symbol.Scope = GlobalScope
	}

	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	// obj, ok := s.store[name]

	// outer := s.Outer
	// for !ok && outer != nil {
	// 	obj, ok = outer.Resolve(name)
	// 	outer = outer.Outer
	// }

	// 递归结构
	obj, ok := s.store[name]
	if !ok && s.Outer != nil { // 这一层找不到就向上找，然后立即返回
		obj, ok = s.Outer.Resolve(name)
		return obj, ok
	}
	return obj, ok
}
