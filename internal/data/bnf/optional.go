package bnf

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

var _ BNF = &Optional{}

type Optional struct {
	BNF
}

func (o Optional) Accept(tokensCh *channel.TransactionChannel[option.Option[token.Token]]) error {
	defer tokensCh.Rollback()

	if err := o.BNF.Accept(tokensCh.StartTx()); errors.Is(err, ErrUnexpectedToken) {
		// Return without committing tx.
		return nil
	} else if err != nil {
		return fmt.Errorf("error in optional: %w", err)
	}

	tokensCh.Commit()

	return nil
}
