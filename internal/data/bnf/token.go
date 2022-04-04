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
	log.Debug("accepting")

	t, err := tokensCh.Read().Unwrap()
	if err != nil {
		log.Debugf("error %v, returning", err)
		return fmt.Errorf("token error: %v", err)
	}

	if tk.ID != t.ID {
		log.Debugf("expected %v, got %v, returning", tk, t.ID)
		return fmt.Errorf("expected %v, got %v: %w", tk, t.ID, ErrUnexpectedToken)
	}

	tokensCh.Commit()

	return nil
}

func (tk Token) String() string {
	return fmt.Sprintf("token %q", tk.ID.String())
}