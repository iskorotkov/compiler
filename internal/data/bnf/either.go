package bnf

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/iskorotkov/compiler/internal/analyzers/syntax_neutralizer"
	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

var _ BNF = &Either{}

type Either struct {
	Name string
	BNFs []BNF
}

func (e Either) Accept(log *zap.SugaredLogger, tokensCh *channel.TransactionChannel[option.Option[token.Token]], neutralizer syntax_neutralizer.Neutralizer) error {
	defer tokensCh.Rollback()

	log = log.Named(e.String())

	var lastError error
	for _, item := range e.BNFs {
		if err := item.Accept(log, tokensCh.StartTx(), neutralizer); errors.Is(err, ErrUnexpectedToken) {
			lastError = err
			log.Infof("%v in %v, skipping", err, e)
			continue
		} else if err != nil {
			log.Warnf("%v in %v, returning", err, e)
			return err
		}

		log.Infof("commit")
		tokensCh.Commit()

		return nil
	}

	return fmt.Errorf("the token is not in a list of expected tokens: %w", lastError)
}

func (e Either) String() string {
	if e.Name != "" {
		return e.Name
	} else {
		return fmt.Sprintf("either")
	}
}
