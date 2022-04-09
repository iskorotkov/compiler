package bnf

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/compiler/internal/contexts"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/modules/syntax_neutralizer"
)

var _ BNF = &Token{}

type Token struct {
	token.ID
}

func (tk Token) Accept(ctx interface {
	contexts.LoggerContext
	contexts.TxChannelContext
	contexts.NeutralizerContext
}) error {
	defer ctx.TxChannel().Rollback()

	ctx, cancel := contexts.Scoped(ctx, tk.String())
	defer cancel()

	t, err := ctx.TxChannel().Read().Unwrap()
	if err != nil {
		ctx.Logger().Warnf("error %v, returning", err)
		return fmt.Errorf("token error: %v", err)
	}

	_, err = ctx.Neutralizer().Neutralize(tk.ID, t)
	if err != nil {
		if errors.Is(err, syntax_neutralizer.UnfixableError) {
			ctx.Logger().Warnf("unfixable syntax error: %v", err)
			return fmt.Errorf("%v: expected %q, got %q: %w", t.Literal, tk, t.ID, ErrUnexpectedToken)
		}

		ctx.Logger().Infof("fixed syntax error: %v", err)
	}

	ctx.Logger().Infof("commit")
	ctx.TxChannel().Commit()

	return nil
}

func (tk Token) String() string {
	return tk.ID.String()
}
