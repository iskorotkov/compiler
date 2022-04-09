package contexts

import (
	"context"

	"go.uber.org/zap"

	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channels"
	"github.com/iskorotkov/compiler/internal/fn/options"
	"github.com/iskorotkov/compiler/internal/modules/syntax_neutralizer"
)

type FullContext interface {
	context.Context
	LoggerContext
	ErrorsContext
}

type LoggerContext interface {
	Logger() *zap.SugaredLogger
	SetLogger(logger *zap.SugaredLogger)
}

type ErrorsContext interface {
	AddError(err error)
	Errors() []error
}

type TxChannelContext interface {
	TxChannel() *channels.TxChannel[options.Option[token.Token]]
}

type NeutralizerContext interface {
	Neutralizer() syntax_neutralizer.Neutralizer
}
