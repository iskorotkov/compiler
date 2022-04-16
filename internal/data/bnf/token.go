package bnf

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channel"
	"github.com/iskorotkov/compiler/internal/fn/option"
	"github.com/iskorotkov/compiler/internal/module/syntax_neutralizer"
)

var _ BNF = Token{}

type Token struct {
	ID token.ID
	ast.Markers
}

func (tk Token) Build(ctx interface {
	context.LoggerContext
	context.NeutralizerContext
}, ch *channel.TxChannel[option.Option[token.Token]]) (ast.Node, error) {
	defer ch.Rollback()

	ctx, cancel := context.Scoped(ctx, tk.String())
	defer cancel()

	t, err := ch.Read().Unwrap()
	if err != nil {
		ctx.Logger().Warnf("error %v, returning", err)
		return nil, fmt.Errorf("token error: %v", err)
	}

	_, err = ctx.Neutralizer().Neutralize(tk.ID, t)
	if err != nil {
		if errors.Is(err, syntax_neutralizer.UnfixableError) {
			ctx.Logger().Warnf("unfixable syntax error: %v", err)
			return nil, fmt.Errorf("%v: expected %q, got %q: %w", t.Literal, tk, t.ID, ErrUnexpectedToken)
		}

		ctx.Logger().Infof("fixed syntax error: %v", err)
	}

	ctx.Logger().Infof("commit")
	ch.Commit()

	return ast.Token(t, tk.Markers), nil
}

func (tk Token) String() string {
	return tk.ID.String()
}
