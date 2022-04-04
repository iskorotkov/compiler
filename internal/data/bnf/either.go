package bnf

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

var _ BNF = &Either{}

type Either struct {
	Name string
	BNFs []BNF
}

func (e Either) Accept(log *zap.SugaredLogger, tokensCh *channel.TransactionChannel[option.Option[token.Token]]) error {
	defer tokensCh.Rollback()

	log = log.Named(e.String())
	log.Debug("accepting")

	var lastError error
	for _, item := range e.BNFs {
		if err := item.Accept(log, tokensCh.StartTx()); errors.Is(err, ErrUnexpectedToken) {
			lastError = err
			log.Debugf("error %v, skipping", err)
			continue
		} else if err != nil {
			log.Debugf("error %v, returning", err)
			return fmt.Errorf("error in optional: %w", err)
		}

		tokensCh.Commit()

		return nil
	}

	return fmt.Errorf("the token is not in a list of expected tokens: %w", lastError)
}

func (e Either) String() string {
	if e.Name != "" {
		return fmt.Sprintf("either %q", e.Name)
	} else {
		return fmt.Sprintf("either")
	}
}
