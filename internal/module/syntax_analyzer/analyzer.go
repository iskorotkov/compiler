package syntax_analyzer

import (
	"errors"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
	"github.com/iskorotkov/compiler/internal/data/bnf"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channel"
	"github.com/iskorotkov/compiler/internal/module/neutralizer"
)

type SyntaxAnalyzer struct {
	buffer      int
	neutralizer neutralizer.Neutralizer
}

func New(buffer int) *SyntaxAnalyzer {
	return &SyntaxAnalyzer{
		buffer: buffer,
	}
}

func (a SyntaxAnalyzer) Analyze(
	ctx interface {
		context.LoggerContext
		context.NeutralizerContext
		context.ErrorsContext
	},
	input <-chan token.Token,
) <-chan ast.Node {
	ch := make(chan ast.Node, a.buffer)

	go func() {
		defer close(ch)

		ctx.Logger().Infof("syntax analysis started")

		res, err := bnf.Program.Build(ctx, channel.NewTxChannel(input))
		if err != nil {
			ctx.Logger().Errorf("error during syntax analysis: %v", err)

			var unexpectedTokenError *bnf.UnexpectedTokenError
			if errors.As(err, &unexpectedTokenError) {
				ctx.AddError(context.ErrorSourceSyntax, unexpectedTokenError.Actual.Position, err)
			}

			return
		}

		ctx.Logger().Infof("syntax analysis succeeded")
		ch <- res
	}()

	return ch
}
