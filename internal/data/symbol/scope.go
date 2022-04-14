package symbol

import (
	"fmt"
)

type Scope struct {
	parent  *Scope
	symbols map[int]Symbol
}

func NewScope() Scope {
	scope := Scope{
		parent:  nil,
		symbols: map[int]Symbol{},
	}

	integerSymbol := Type{Token: builtinToken("integer")}
	realSymbol := Type{Token: builtinToken("real")}
	booleanSymbol := Type{Token: builtinToken("boolean")}
	stringSymbol := Type{Token: builtinToken("string")}
	voidSymbol := Type{Token: builtinToken("void")}
	writelnSymbol := Func{
		Token: builtinToken("writeln"),
		Params: []Var{
			{Token: builtinToken("s"), Type: stringSymbol},
		},
		ReturnType: voidSymbol,
	}

	_ = scope.Add(&integerSymbol)
	_ = scope.Add(&realSymbol)
	_ = scope.Add(&booleanSymbol)
	_ = scope.Add(&voidSymbol)
	_ = scope.Add(&writelnSymbol)

	return scope
}

func (s Scope) Lookup(symbol Symbol) (Symbol, bool) {
	if symbol, ok := s.symbols[symbol.Hash()]; ok {
		return symbol, true
	}

	if s.parent != nil {
		return s.parent.Lookup(symbol)
	}

	return nil, false
}

func (s Scope) Add(symbol Symbol) error {
	if _, ok := s.symbols[symbol.Hash()]; ok {
		return fmt.Errorf("%v was already declared in this scope", symbol)
	}

	s.symbols[symbol.Hash()] = symbol

	return nil
}

func (s *Scope) SubScope(symbols []Symbol) *Scope {
	m := make(map[int]Symbol)
	for _, symbol := range symbols {
		m[symbol.Hash()] = symbol
	}

	return &Scope{
		parent:  s,
		symbols: m,
	}
}

func (s Scope) ParentScope() *Scope {
	return s.parent
}

func (s Scope) Symbols() []Symbol {
	var symbols []Symbol
	for _, symbol := range s.symbols {
		symbols = append(symbols, symbol)
	}

	return symbols
}
