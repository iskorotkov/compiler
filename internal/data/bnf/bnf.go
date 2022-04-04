package bnf

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

type BNF interface {
	fmt.Stringer
	Accept(log *zap.SugaredLogger, tokensCh *channel.TransactionChannel[option.Option[token.Token]]) error
}
