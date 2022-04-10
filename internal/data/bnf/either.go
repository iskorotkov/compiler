package bnf

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/compiler/internal/context"
)

var _ BNF = &Either{}

type Either struct {
	Name string
	BNFs []BNF
}

func (e Either) Accept(ctx interface {
	context.LoggerContext
	context.TxChannelContext
	context.NeutralizerContext
}) error {
	defer ctx.TxChannel().Rollback()

	ctx, cancel := context.Scoped(ctx, e.Name)
	defer cancel()

	var lastError error
	for _, item := range e.BNFs {
		if err := item.Accept(ctx); errors.Is(err, ErrUnexpectedToken) {
			lastError = err
			ctx.Logger().Infof("%v in %v, skipping", err, e)
			continue
		} else if err != nil {
			ctx.Logger().Warnf("%v in %v, returning", err, e)
			return err
		}

		ctx.Logger().Infof("commit")
		ctx.TxChannel().Commit()

		return nil
	}

	return fmt.Errorf("the token is not in a list of expected tokens: %w", lastError)
}

func (e Either) String() string {
	if e.Name != "" {
		return e.Name
	} else {
		return fmt.Sprintf("either")
	}
}
