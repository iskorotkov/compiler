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

var _ BNF = &Several{}

type Several struct {
	Name string
	BNF
}

func (s Several) Accept(log *zap.SugaredLogger, tokensCh *channel.TransactionChannel[option.Option[token.Token]], neutralizer syntax_neutralizer.Neutralizer) error {
	defer tokensCh.Rollback()

	log = log.Named(s.String())

	for {
		if err := s.BNF.Accept(log, tokensCh.StartTx(), neutralizer); errors.Is(err, ErrUnexpectedToken) {
			tokensCh.Commit()
			log.Infof("%v in %v, committing tx", err, s)
			return nil
		} else if err != nil {
			log.Warnf("%v in %v, returning", err, s)
			return err
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
