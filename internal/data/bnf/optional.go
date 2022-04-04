package bnf

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

var _ BNF = &Optional{}

type Optional struct {
	Name string
	BNF
}

func (o Optional) Accept(log *zap.SugaredLogger, tokensCh *channel.TransactionChannel[option.Option[token.Token]]) error {
	defer tokensCh.Rollback()

	log = log.Named(o.String())

	if err := o.BNF.Accept(log, tokensCh.StartTx()); errors.Is(err, ErrUnexpectedToken) {
		log.Debugf("%v in %s, rollback tx", o, err)
		return nil
	} else if err != nil {
		log.Debugf("%v, returning", err)
		return fmt.Errorf("error in %s: %w", o, err)
	}

	log.Debugf("commit")
	tokensCh.Commit()

	return nil
}

func (o Optional) String() string {
	if o.Name != "" {
		return o.Name
	} else {
		return fmt.Sprintf("optional")
	}
}
