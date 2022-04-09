package syntax_analyzer

import (
	"github.com/iskorotkov/compiler/internal/analyzers/syntax_neutralizer"
	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/bnf"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
	"github.com/iskorotkov/compiler/internal/logger"
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

func (a SyntaxAnalyzer) Analyze(input <-chan option.Option[token.Token]) <-chan option.Option[bnf.BNF] {
	ch := make(chan option.Option[bnf.BNF], a.buffer)

	go func() {
		defer close(ch)

		log.Infof("syntax analysis started")

		tx := channel.NewTransactionChannel(input)

		if err := bnf.Program.Accept(log, tx, a.neutralizer); err != nil {
			log.Errorf("error during syntax analysis: %v", err)
			ch <- option.Err[bnf.BNF](err)
			return
		}

		log.Infof("syntax analysis succeeded")
		ch <- option.Ok[bnf.BNF](bnf.Program)
	}()

	return ch
}
