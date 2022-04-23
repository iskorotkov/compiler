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

var _ BNF = Either{}

type Either struct {
	Name string
	BNFs []BNF
	ast.Markers
}

func (e Either) Build(ctx interface {
	context.LoggerContext
	context.NeutralizerContext
}, ch *channel.TxChannel[option.Option[token.Token]]) (ast.Node, error) {
	ctx, cancel := context.Scoped(ctx, e.Name)
	defer cancel()

	var lastError error
	for _, item := range e.BNFs {
		ch := ch.StartTx()

		res, err := item.Build(ctx, ch)
		if err != nil {
			ch.Rollback()

			if errors.Is(err, ErrUnexpectedToken) {
				lastError = err
				ctx.Logger().Debugf("%v in %v, skipping", err, e)
				continue
			} else {
				ctx.Logger().Debugf("%v in %v, returning", err, e)
				return nil, err
			}
		}

		ctx.Logger().Debugf("ok, commit tx")
		ch.Commit()

		if res == nil {
			return nil, nil
		}

		return ast.Wrap(res, e.Markers), nil
	}

	ctx.Logger().Errorf("the token is not in a list of expected tokens: %w", lastError)
	return nil, fmt.Errorf("the token is not in a list of expected tokens: %w", lastError)
}

func (e Either) String() string {
	if e.Name != "" {
		return e.Name
	} else {
		return fmt.Sprintf("either")
	}
}
