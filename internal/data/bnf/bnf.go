package bnf

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channels"
	"github.com/iskorotkov/compiler/internal/fn/options"
	"github.com/iskorotkov/compiler/internal/modules/syntax_neutralizer"
)

type BNF interface {
	fmt.Stringer
	Accept(log *zap.SugaredLogger, tokensCh *channels.TxChannel[options.Option[token.Token]], neutralizer syntax_neutralizer.Neutralizer) error
}
