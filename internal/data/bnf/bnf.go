package bnf

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
)

type BNF interface {
	Build(ctx interface {
		context.LoggerContext
		context.TxChannelContext
		context.NeutralizerContext
	}) (ast.Node, error)
	fmt.Stringer
}
