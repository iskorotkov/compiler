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
	defer ch.Rollback()

	ctx, cancel := context.Scoped(ctx, o.Name)
	defer cancel()

	res, err := o.BNF.Build(ctx, ch)
	if errors.Is(err, ErrUnexpectedToken) {
		ctx.Logger().Infof("%v in %v, rollback tx", err, o)
		return nil, nil
	} else if err != nil {
		ctx.Logger().Warnf("%v in %v, returning", err, o)
		return nil, err
	}

	ctx.Logger().Infof("commit")
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
