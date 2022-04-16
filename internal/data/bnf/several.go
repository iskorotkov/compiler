package bnf

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channel"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

var _ BNF = Several{}

type Several struct {
	Name string
	BNF  BNF
	ast.Markers
}

func (s Several) Build(ctx interface {
	context.LoggerContext
	context.NeutralizerContext
}, ch *channel.TxChannel[option.Option[token.Token]]) (ast.Node, error) {
	defer ch.Rollback()

	ctx, cancel := context.Scoped(ctx, s.Name)
	defer cancel()

	var items []ast.Node
	for {
		res, err := s.BNF.Build(ctx, ch)
		if errors.Is(err, ErrUnexpectedToken) {
			ch.Commit()
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
