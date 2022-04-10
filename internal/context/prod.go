package context

import (
	"context"

	"go.uber.org/zap"
)

var _ FullContext = (*prodContext)(nil)

type prodContext struct {
	context.Context
	logger *zap.SugaredLogger
	errors []error
}

func (p *prodContext) setLogger(logger *zap.SugaredLogger) {
	p.logger = logger
}

func (p *prodContext) Logger() *zap.SugaredLogger {
	return p.logger
}

func (p *prodContext) AddError(err error) {
	p.errors = append(p.errors, err)
}

func (p *prodContext) Errors() []error {
	return p.errors
}

func NewProdContext(ctx context.Context) FullContext {
	return &prodContext{
		Context: ctx,
		logger:  zap.NewNop().Sugar(),
		errors:  nil,
	}
}
