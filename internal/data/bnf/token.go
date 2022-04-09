package bnf

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channels"
	"github.com/iskorotkov/compiler/internal/fn/options"
	"github.com/iskorotkov/compiler/internal/modules/syntax_neutralizer"
)

var _ BNF = &Token{}

type Token struct {
	token.ID
}

func (tk Token) Accept(log *zap.SugaredLogger, tokensCh *channels.TxChannel[options.Option[token.Token]], neutralizer syntax_neutralizer.Neutralizer) error {
	defer tokensCh.Rollback()

	log = log.Named(tk.String())

	t, err := tokensCh.Read().Unwrap()
	if err != nil {
		log.Warnf("error %v, returning", err)
		return fmt.Errorf("token error: %v", err)
	}

	_, err = neutralizer.Neutralize(tk.ID, t)
	if err != nil {
		if errors.Is(err, syntax_neutralizer.UnfixableError) {
			log.Warnf("unfixable syntax error: %v", err)
			return fmt.Errorf("%v: expected %q, got %q: %w", t.Literal, tk, t.ID, ErrUnexpectedToken)
		}

		log.Infof("fixed syntax error: %v", err)
	}

	log.Infof("commit")
	tokensCh.Commit()

	return nil
}

func (tk Token) String() string {
	return tk.ID.String()
}
