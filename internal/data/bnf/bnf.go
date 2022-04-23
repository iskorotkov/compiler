package bnf

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channel"
)

type BNF interface {
	Build(ctx interface {
		context.LoggerContext
		context.NeutralizerContext
	}, ch *channel.TxChannel[token.Token]) (ast.Node, error)
	fmt.Stringer
}
