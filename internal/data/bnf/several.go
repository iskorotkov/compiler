package bnf

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
)

var _ BNF = Several{}

type Several struct {
	Name string
	BNF  BNF
	ast.Markers
}

func (s Several) Build(ctx interface {
	context.LoggerContext
	context.TxChannelContext
	context.NeutralizerContext
}) (ast.Node, error) {
	defer ctx.TxChannel().Rollback()

	ctx, cancel := context.Scoped(ctx, s.Name)
	defer cancel()

	var items []ast.Node
	for {
		res, err := s.BNF.Build(ctx)
		if errors.Is(err, ErrUnexpectedToken) {
			ctx.TxChannel().Commit()
			ctx.Logger().Infof("%v in %v, committing tx", err, s)

			return ast.WrapSlice(items, s.Markers), nil
		} else if err != nil {
			ctx.Logger().Warnf("%v in %v, returning", err, s)
			return nil, err
		}

		if res != nil {
			items = append(items, res)
		}
	}
}

func (s Several) String() string {
	if s.Name != "" {
		return s.Name
	} else {
		return fmt.Sprintf("several")
	}
}
