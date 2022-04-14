package bnf

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
)

var _ BNF = Sequence{}

type Sequence struct {
	Name string
	BNFs []BNF
	ast.Markers
}

func (s Sequence) Build(ctx interface {
	context.LoggerContext
	context.TxChannelContext
	context.NeutralizerContext
}) (ast.Node, error) {
	defer ctx.TxChannel().Rollback()

	ctx, cancel := context.Scoped(ctx, s.Name)
	defer cancel()

	var items []ast.Node
	for _, item := range s.BNFs {
		res, err := item.Build(ctx)
		if err != nil {
			ctx.Logger().Warnf("%v in %v, returning", err, s)
			return nil, err
		}

		if res != nil {
			items = append(items, res)
		}
	}

	ctx.Logger().Infof("commit")
	ctx.TxChannel().Commit()

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
