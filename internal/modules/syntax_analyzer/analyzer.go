package syntax_analyzer

import (
	"github.com/iskorotkov/compiler/internal/data/bnf"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channels"
	"github.com/iskorotkov/compiler/internal/fn/options"
	"github.com/iskorotkov/compiler/internal/logger"
	"github.com/iskorotkov/compiler/internal/modules/syntax_neutralizer"
)

var log = logger.New().Named("syntax_analyzer")

type SyntaxAnalyzer struct {
	buffer      int
	neutralizer syntax_neutralizer.Neutralizer
}

func New(buffer int, neutralizationMaxDistance int) *SyntaxAnalyzer {
	return &SyntaxAnalyzer{
		buffer:      buffer,
		neutralizer: syntax_neutralizer.New(neutralizationMaxDistance),
	}
}

func (a SyntaxAnalyzer) Analyze(input <-chan options.Option[token.Token]) <-chan options.Option[bnf.BNF] {
	ch := make(chan options.Option[bnf.BNF], a.buffer)

	go func() {
		defer close(ch)

		log.Infof("syntax analysis started")

		tx := channels.NewTxChannel(input)

		if err := bnf.Program.Accept(log, tx, a.neutralizer); err != nil {
			log.Errorf("error during syntax analysis: %v", err)
			ch <- options.Err[bnf.BNF](err)
			return
		}

		log.Infof("syntax analysis succeeded")
		ch <- options.Ok[bnf.BNF](bnf.Program)
	}()

	return ch
}
