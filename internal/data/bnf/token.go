package bnf

import (
	"errors"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channel"
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
}, ch *channel.TxChannel[token.Token]) (ast.Node, error) {
	ctx, cancel := context.Scoped(ctx, tk.String())
	defer cancel()

	t := ch.Read()
	if _, err := ctx.Neutralizer().Neutralize(tk.ID, t); err != nil {
		if errors.Is(err, syntax_neutralizer.UnfixableError) {
			ctx.Logger().Warnf("unfixable syntax error: %v", err)
			return nil, &UnexpectedTokenError{
				Expected: tk.ID,
				Actual:   t,
			}
		}

		ctx.Logger().Infof("fixed syntax error: %v", err)
	}

	ctx.Logger().Infof("ok")

	return ast.Token(t, tk.Markers), nil
}

func (tk Token) String() string {
	return tk.ID.String()
}
