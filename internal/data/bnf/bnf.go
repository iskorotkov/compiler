package bnf

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channel"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

type BNF interface {
	Build(ctx interface {
		context.LoggerContext
		context.NeutralizerContext
	}, ch *channel.TxChannel[option.Option[token.Token]]) (ast.Node, error)
	fmt.Stringer
}
