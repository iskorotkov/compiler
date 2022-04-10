package context

import (
	"context"

	"go.uber.org/zap"

	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channel"
	"github.com/iskorotkov/compiler/internal/fn/option"
	"github.com/iskorotkov/compiler/internal/module/syntax_neutralizer"
)

type FullContext interface {
	context.Context
	LoggerContext
	ErrorsContext
}

type LoggerContext interface {
	Logger() *zap.SugaredLogger
	setLogger(logger *zap.SugaredLogger)
}

type ErrorsContext interface {
	AddError(err error)
	Errors() []error
}

type TxChannelContext interface {
	TxChannel() *channel.TxChannel[option.Option[token.Token]]
}

type NeutralizerContext interface {
	Neutralizer() syntax_neutralizer.Neutralizer
}
