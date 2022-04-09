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

var _ BNF = &Optional{}

type Optional struct {
	Name string
	BNF
}

func (o Optional) Accept(log *zap.SugaredLogger, tokensCh *channels.TxChannel[options.Option[token.Token]], neutralizer syntax_neutralizer.Neutralizer) error {
	defer tokensCh.Rollback()

	log = log.Named(o.String())

	if err := o.BNF.Accept(log, tokensCh.StartTx(), neutralizer); errors.Is(err, ErrUnexpectedToken) {
		log.Infof("%v in %v, rollback tx", err, o)
		return nil
	} else if err != nil {
		log.Warnf("%v in %v, returning", err, o)
		return err
	}

	log.Infof("commit")
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
