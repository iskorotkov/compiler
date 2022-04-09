package bnf

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/contexts"
)

var _ BNF = &Sequence{}

type Sequence struct {
	Name string
	BNFs []BNF
}

func (s Sequence) Accept(ctx interface {
	contexts.LoggerContext
	contexts.TxChannelContext
	contexts.NeutralizerContext
}) error {
	defer ctx.TxChannel().Rollback()

	ctx, cancel := contexts.Scoped(ctx, s.String())
	defer cancel()

	for _, item := range s.BNFs {
		if err := item.Accept(ctx); err != nil {
			ctx.Logger().Warnf("%v in %v, returning", err, s)
			return err
		}
	}

	ctx.Logger().Infof("commit")
	ctx.TxChannel().Commit()

	return nil
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
