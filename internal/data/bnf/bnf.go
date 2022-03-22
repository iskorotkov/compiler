package bnf

import (
	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

type BNF interface {
	Accept(tokensCh *channel.TransactionChannel[option.Option[token.Token]]) error
}
