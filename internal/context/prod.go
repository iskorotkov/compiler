package context

import (
	"context"

	"go.uber.org/zap"

	"github.com/iskorotkov/compiler/internal/data/symbol"
	"github.com/iskorotkov/compiler/internal/module/syntax_neutralizer"
)

var _ FullContext = (*prodContext)(nil)

type prodContext struct {
	context.Context
	errorsContext
	neutralizerContext
	symbolScopeContext
	logger *zap.SugaredLogger
}

func (c *prodContext) setLogger(logger *zap.SugaredLogger) {
	c.logger = logger
}

func (c *prodContext) Logger() *zap.SugaredLogger {
	return c.logger
}

func NewProdContext(ctx context.Context) FullContext {
	return &prodContext{
		Context:            ctx,
		logger:             zap.NewNop().Sugar(),
		errorsContext:      errorsContext{},
		neutralizerContext: neutralizerContext{neutralizer: syntax_neutralizer.New(1)},
		symbolScopeContext: symbolScopeContext{scope: symbol.NewScope()},
	}
}
