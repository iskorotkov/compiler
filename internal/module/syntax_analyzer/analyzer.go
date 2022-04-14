package syntax_analyzer

import (
	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
	"github.com/iskorotkov/compiler/internal/data/bnf"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channel"
	"github.com/iskorotkov/compiler/internal/fn/option"
	"github.com/iskorotkov/compiler/internal/module/syntax_neutralizer"
)

type SyntaxAnalyzer struct {
	buffer      int
	neutralizer syntax_neutralizer.Neutralizer
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
	},
	input <-chan option.Option[token.Token],
) <-chan option.Option[ast.Node] {
	ch := make(chan option.Option[ast.Node], a.buffer)

	go func() {
		defer close(ch)

		ctx.Logger().Infof("syntax analysis started")

		res, err := bnf.Program.Build(struct {
			context.LoggerContext
			context.NeutralizerContext
			context.TxChannelContext
		}{
			LoggerContext:      ctx,
			NeutralizerContext: ctx,
			TxChannelContext:   context.NewTxChannelContext(channel.NewTxChannel(input)),
		})

		if err != nil {
			ctx.Logger().Errorf("error during syntax analysis: %v", err)
			ch <- option.Err[ast.Node](err)
			return
		}

		ctx.Logger().Infof("syntax analysis succeeded")
		ch <- option.Ok[ast.Node](res)
	}()

	return ch
}
