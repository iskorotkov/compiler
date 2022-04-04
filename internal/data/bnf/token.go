package bnf

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

var _ BNF = &Token{}

type Token struct {
	token.ID
}

func (tk Token) Accept(log *zap.SugaredLogger, tokensCh *channel.TransactionChannel[option.Option[token.Token]]) error {
	defer tokensCh.Rollback()

	log = log.Named(tk.String())

	t, err := tokensCh.Read().Unwrap()
	if err != nil {
		log.Debugf("error %v, returning", err)
		return fmt.Errorf("token error: %v", err)
	}

	if tk.ID != t.ID {
		log.Debugf("expected %v, got %v, returning", tk, t.ID)
		return fmt.Errorf("%v: expected %q, got %q: %w", t.Literal, tk, t.ID, ErrUnexpectedToken)
	}

	log.Debugf("commit")
	tokensCh.Commit()

	return nil
}

func (tk Token) String() string {
	return tk.ID.String()
}
