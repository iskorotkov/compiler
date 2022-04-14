package context

import (
	"github.com/iskorotkov/compiler/internal/data/symbol"
)

var _ SymbolScopeContext = (*symbolScopeContext)(nil)

type symbolScopeContext struct {
	scope symbol.Scope
}

func (t *symbolScopeContext) SymbolScope() symbol.Scope {
	return t.scope
}
