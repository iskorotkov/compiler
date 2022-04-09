package bnf

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/compiler/internal/contexts"
)

var _ BNF = &Either{}

type Either struct {
	Name string
	BNFs []BNF
}

func (e Either) Accept(ctx interface {
	contexts.LoggerContext
	contexts.TxChannelContext
	contexts.NeutralizerContext
}) error {
	defer ctx.TxChannel().Rollback()

	ctx, cancel := contexts.Scoped(ctx, e.String())
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
