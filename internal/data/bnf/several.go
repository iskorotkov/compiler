package bnf

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/compiler/internal/contexts"
)

var _ BNF = &Several{}

type Several struct {
	Name string
	BNF
}

func (s Several) Accept(ctx interface {
	contexts.LoggerContext
	contexts.TxChannelContext
	contexts.NeutralizerContext
}) error {
	defer ctx.TxChannel().Rollback()

	ctx, cancel := contexts.Scoped(ctx, s.String())
	defer cancel()

	for {
		if err := s.BNF.Accept(ctx); errors.Is(err, ErrUnexpectedToken) {
			ctx.TxChannel().Commit()
			ctx.Logger().Infof("%v in %v, committing tx", err, s)
			return nil
		} else if err != nil {
			ctx.Logger().Warnf("%v in %v, returning", err, s)
			return err
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
