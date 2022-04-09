package bnf

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channels"
	"github.com/iskorotkov/compiler/internal/fn/options"
	"github.com/iskorotkov/compiler/internal/modules/syntax_neutralizer"
)

var _ BNF = &Sequence{}

type Sequence struct {
	Name string
	BNFs []BNF
}

func (s Sequence) Accept(log *zap.SugaredLogger, tokensCh *channels.TxChannel[options.Option[token.Token]], neutralizer syntax_neutralizer.Neutralizer) error {
	defer tokensCh.Rollback()

	log = log.Named(s.String())

	for _, item := range s.BNFs {
		if err := item.Accept(log, tokensCh.StartTx(), neutralizer); err != nil {
			log.Warnf("%v in %v, returning", err, s)
			return err
		}
	}

	log.Infof("commit")
	tokensCh.Commit()

	return nil
}

func (s Sequence) String() string {
	if len(s.BNFs) == 0 {
		return "empty"
	}

	if s.Name != "" {
		return s.Name
	} else {
		return fmt.Sprintf("sequence")
	}
}
