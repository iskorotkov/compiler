package bnf

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

var _ BNF = &Token{0}

type Token struct {
	token.ID
}

func (tk Token) Accept(tokensCh *channel.TransactionChannel[option.Option[token.Token]]) error {
	defer tokensCh.Rollback()

	opt := tokensCh.Read()
	t, err := opt.Unwrap()
	if err != nil {
		return fmt.Errorf("token error: %v", err)
	}

	if tk.ID != t.ID {
		return fmt.Errorf("expected token %v, got %v: %w", tk, t.ID, ErrUnexpectedToken)
	}

	tokensCh.Commit()

	return nil
}
