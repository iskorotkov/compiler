package bnf

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
)

var _ BNF = Either{}

type Either struct {
	Name string
	BNFs []BNF
	ast.Markers
}

func (e Either) Build(ctx interface {
	context.LoggerContext
	context.TxChannelContext
	context.NeutralizerContext
}) (ast.Node, error) {
	defer ctx.TxChannel().Rollback()

	ctx, cancel := context.Scoped(ctx, e.Name)
	defer cancel()

	var lastError error
	for _, item := range e.BNFs {
		res, err := item.Build(ctx)
		if errors.Is(err, ErrUnexpectedToken) {
			lastError = err
			ctx.Logger().Infof("%v in %v, skipping", err, e)
			continue
		} else if err != nil {
			ctx.Logger().Warnf("%v in %v, returning", err, e)
			return nil, err
		}

		ctx.Logger().Infof("commit")
		ctx.TxChannel().Commit()

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
