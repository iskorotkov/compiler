package aggregates

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/bnf"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

var _ bnf.BNF = &Token{}

type Token struct {
	token.ID
}

func (tk Token) Accept(tokensCh *channel.TransactionChannel[option.Option[token.Token]]) error {
	defer tokensCh.Rollback()

	log.Printf("expecting %v", tk)

	opt := tokensCh.Read()
	t, err := opt.Unwrap()
	if err != nil {
		return fmt.Errorf("token error: %v", err)
	}

	if tk.ID != t.ID {
		return fmt.Errorf("expected token %v, got %v: %w", tk, t.ID, bnf.ErrUnexpectedToken)
	}

	tokensCh.Commit()

	return nil
}

func (tk Token) String() string {
	return tk.ID.String()
}
