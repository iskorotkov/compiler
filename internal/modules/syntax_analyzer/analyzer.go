package syntax_analyzer

import (
	"github.com/iskorotkov/compiler/internal/contexts"
	"github.com/iskorotkov/compiler/internal/data/bnf"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channels"
	"github.com/iskorotkov/compiler/internal/fn/options"
	"github.com/iskorotkov/compiler/internal/modules/syntax_neutralizer"
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
		contexts.LoggerContext
		contexts.NeutralizerContext
	},
	input <-chan options.Option[token.Token],
) <-chan options.Option[bnf.BNF] {
	ch := make(chan options.Option[bnf.BNF], a.buffer)

	go func() {
		defer close(ch)

		ctx.Logger().Infof("syntax analysis started")

		tx := channels.NewTxChannel(input)

		if err := bnf.Program.Accept(struct {
			contexts.LoggerContext
			contexts.NeutralizerContext
			contexts.TxChannelContext
		}{ctx, ctx, contexts.NewTxChannelContext(tx)}); err != nil {
			ctx.Logger().Errorf("error during syntax analysis: %v", err)
			ch <- options.Err[bnf.BNF](err)
			return
		}

		ctx.Logger().Infof("syntax analysis succeeded")
		ch <- options.Ok[bnf.BNF](bnf.Program)
	}()

	return ch
}
