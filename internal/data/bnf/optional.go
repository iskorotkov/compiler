package bnf

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/compiler/internal/contexts"
)

var _ BNF = &Optional{}

type Optional struct {
	Name string
	BNF
}

func (o Optional) Accept(ctx interface {
	contexts.LoggerContext
	contexts.TxChannelContext
	contexts.NeutralizerContext
}) error {
	defer ctx.TxChannel().Rollback()

	ctx, cancel := contexts.Scoped(ctx, o.String())
	defer cancel()

	if err := o.BNF.Accept(ctx); errors.Is(err, ErrUnexpectedToken) {
		ctx.Logger().Infof("%v in %v, rollback tx", err, o)
		return nil
	} else if err != nil {
		ctx.Logger().Warnf("%v in %v, returning", err, o)
		return err
	}

	ctx.Logger().Infof("commit")
	ctx.TxChannel().Commit()

	return nil
}

func (o Optional) String() string {
	if o.Name != "" {
		return o.Name
	} else {
		return fmt.Sprintf("optional")
	}
}
