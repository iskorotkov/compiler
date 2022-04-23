package bnf

import (
	"errors"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channel"
	"github.com/iskorotkov/compiler/internal/module/neutralizer"
)

var _ BNF = Token{}

type Token struct {
	ID token.ID
	ast.Markers
}

func (tk Token) Build(ctx interface {
	context.LoggerContext
	context.NeutralizerContext
	context.ErrorsContext
}, ch *channel.TxChannel[token.Token]) (ast.Node, error) {
	ctx, cancel := context.Scoped(ctx, tk.String())
	defer cancel()

	actualToken := ch.Read()

	fixedToken, err := ctx.Neutralizer().NeutralizeKeyword(tk.ID, actualToken)
	if err != nil {
		if errors.Is(err, &neutralizer.UnfixableKeywordError{}) {
			ctx.Logger().Warnf("unfixable syntax error: %v", err)
			return nil, &UnexpectedTokenError{
				Expected: tk.ID,
				Actual:   actualToken,
			}
		}

		ctx.Logger().Infof("fixed syntax error: %v", err)
		ctx.AddError(context.ErrorSourceSyntax, actualToken.Position, err)
	}

	ctx.Logger().Infof("ok")

	return ast.Token(fixedToken, tk.Markers), nil
}

func (tk Token) String() string {
	return tk.ID.String()
}
