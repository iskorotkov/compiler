package syntax_analyzer

import (
	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/bnf"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
	"github.com/iskorotkov/compiler/internal/logger"
)

var log = logger.New("syntax analyzer")

type SyntaxAnalyzer struct {
	buffer int
}

func New(buffer int) *SyntaxAnalyzer {
	return &SyntaxAnalyzer{
		buffer: buffer,
	}
}

func (a SyntaxAnalyzer) Analyze(input <-chan option.Option[token.Token]) <-chan option.Option[bnf.BNF] {
	ch := make(chan option.Option[bnf.BNF], a.buffer)

	go func() {
		defer close(ch)

		log.Printf("started")

		tx := channel.NewTransactionChannel(input)

		if err := bnf.Program.Accept(tx); err != nil {
			log.Print(err)
			ch <- option.Err[bnf.BNF](err)
		}

		log.Printf("success")
		ch <- option.Ok[bnf.BNF](bnf.Program)
	}()

	return ch
}
