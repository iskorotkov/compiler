package bnf

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/contexts"
)

type BNF interface {
	fmt.Stringer
	Accept(
		ctx interface {
			contexts.LoggerContext
			contexts.TxChannelContext
			contexts.NeutralizerContext
		},
	) error
}
