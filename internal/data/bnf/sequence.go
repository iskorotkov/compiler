package bnf

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channel"
)

var _ BNF = Sequence{}

type Sequence struct {
	Name string
	BNFs []BNF
	ast.Markers
}

func (s Sequence) Build(ctx interface {
	context.LoggerContext
	context.NeutralizerContext
}, ch *channel.TxChannel[token.Token]) (ast.Node, error) {
	ctx, cancel := context.Scoped(ctx, s.Name)
	defer cancel()

	var items []ast.Node
	for _, item := range s.BNFs {
		res, err := item.Build(ctx, ch)
		if err != nil {
			ctx.Logger().Debugf("%v in %v, returning", err, s)
			return nil, err
		}

		if res != nil {
			items = append(items, res)
		}
	}

	ctx.Logger().Debugf("ok")

	return ast.WrapSlice(items, s.Markers), nil
}

func (s Sequence) String() string {
	if len(s.BNFs) == 0 {
		return "empty"
	}

	if s.Name != "" {
		return s.Name
	} else {
		return fmt.Sprintf("sequence")
	}
}
