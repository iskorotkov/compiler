package bnf

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channel"
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
}, ch *channel.TxChannel[token.Token]) (ast.Node, error) {
	ctx, cancel := context.Scoped(ctx, s.Name)
	defer cancel()

	var items []ast.Node
	for {
		ch := ch.StartTx()

		res, err := s.BNF.Build(ctx, ch)
		if err != nil {
			ch.Rollback()

			if errors.Is(err, &UnexpectedTokenError{}) {
				ctx.Logger().Debugf("%v in %v, commit tx", err, s)
				return ast.WrapSlice(items, s.Markers), nil
			} else {
				ctx.Logger().Debugf("%v in %v, returning", err, s)
				return nil, err
			}
		}

		ch.Commit()

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
