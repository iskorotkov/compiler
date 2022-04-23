package context

import (
	"context"

	"go.uber.org/zap"

	"github.com/iskorotkov/compiler/internal/data/literal"
	"github.com/iskorotkov/compiler/internal/module/syntax_neutralizer"
)

type FullContext interface {
	context.Context
	LoggerContext
	ErrorsContext
	NeutralizerContext
}

type LoggerContext interface {
	Logger() *zap.SugaredLogger
	setLogger(logger *zap.SugaredLogger)
}

type ErrorsContext interface {
	AddError(source ErrorSource, position literal.Position, err error)
	Errors() []Error
}

type NeutralizerContext interface {
	Neutralizer() syntax_neutralizer.Neutralizer
}
