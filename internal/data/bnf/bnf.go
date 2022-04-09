package bnf

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/context"
)

type BNF interface {
	fmt.Stringer
	Accept(
		ctx interface {
			context.LoggerContext
			context.TxChannelContext
			context.NeutralizerContext
		},
	) error
}
