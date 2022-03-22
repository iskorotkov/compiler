package bnf

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

var _ BNF = &Composite{}

type Composite struct {
	BNFs []BNF
}

func (c Composite) Accept(tokensCh *channel.TransactionChannel[option.Option[token.Token]]) error {
	defer tokensCh.Rollback()

	for i, bnf := range c.BNFs {
		if err := bnf.Accept(tokensCh.StartTx()); err != nil {
			return fmt.Errorf("error in composite at index %d: %w", i, err)
		}
	}

	tokensCh.Commit()

	return nil
}
