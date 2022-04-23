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

var _ BNF = Optional{}

type Optional struct {
	Name string
	BNF  BNF
	ast.Markers
}

func (o Optional) Build(ctx interface {
	context.LoggerContext
	context.NeutralizerContext
}, ch *channel.TxChannel[option.Option[token.Token]]) (ast.Node, error) {
	ctx, cancel := context.Scoped(ctx, o.Name)
	defer cancel()

	ch = ch.StartTx()

	res, err := o.BNF.Build(ctx, ch)
	if err != nil {
		ch.Rollback()

		if errors.Is(err, ErrUnexpectedToken) {
			ctx.Logger().Debugf("%v in %v, rollback tx", err, o)
			return nil, nil
		} else {
			ctx.Logger().Debugf("%v in %v, commit tx", err, o)
			ch.Commit()
			return nil, err
		}
	}

	ctx.Logger().Debugf("ok, commit tx")
	ch.Commit()

	return ast.Wrap(res, o.Markers), nil
}

func (o Optional) String() string {
	if o.Name != "" {
		return o.Name
	} else {
		return fmt.Sprintf("optional")
	}
}
