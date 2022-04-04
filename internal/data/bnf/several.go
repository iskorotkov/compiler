package bnf

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

var _ BNF = &Several{}

type Several struct {
	Name string
	BNF
}

func (s Several) Accept(log *zap.SugaredLogger, tokensCh *channel.TransactionChannel[option.Option[token.Token]]) error {
	defer tokensCh.Rollback()

	log = log.Named(s.String())

	for {
		if err := s.BNF.Accept(log, tokensCh.StartTx()); errors.Is(err, ErrUnexpectedToken) {
			tokensCh.Commit()
			log.Debugf("%v in %v, committing tx", err, s)
			return nil
		} else if err != nil {
			log.Debugf("%v in %v, returning", err, s)
			return fmt.Errorf("error in %v: %w", s, err)
		}
	}
}

func (s Several) String() string {
	if s.Name != "" {
		return s.Name
	} else {
		return fmt.Sprintf("several")
	}
}
