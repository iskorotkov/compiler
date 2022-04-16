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
	defer ch.Rollback()

	ctx, cancel := context.Scoped(ctx, e.Name)
	defer cancel()

	var lastError error
	for _, item := range e.BNFs {
		res, err := item.Build(ctx, ch)
		if errors.Is(err, ErrUnexpectedToken) {
			lastError = err
			ctx.Logger().Infof("%v in %v, skipping", err, e)
			continue
		} else if err != nil {
			ctx.Logger().Warnf("%v in %v, returning", err, e)
			return nil, err
		}

		ctx.Logger().Infof("commit")
		ch.Commit()

		if res == nil {
			return nil, nil
		}

		return ast.Wrap(res, e.Markers), nil
	}

	return nil, fmt.Errorf("the token is not in a list of expected tokens: %w", lastError)
}

func (e Either) String() string {
	if e.Name != "" {
		return e.Name
	} else {
		return fmt.Sprintf("either")
	}
}
