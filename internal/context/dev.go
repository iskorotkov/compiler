package context

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/iskorotkov/compiler/internal/module/syntax_neutralizer"
)

var _ FullContext = (*devContext)(nil)

type devContext struct {
	context.Context
	errorsContext
	neutralizerContext
	logger *zap.SugaredLogger
}

func (c *devContext) setLogger(logger *zap.SugaredLogger) {
	c.logger = logger
}

func (c *devContext) Logger() *zap.SugaredLogger {
	return c.logger
}

func NewDevContext(ctx context.Context) FullContext {
	return &devContext{
		Context: ctx,
		logger: func() *zap.SugaredLogger {
			logger, err := zap.Config{
				Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
				Development:       true,
				DisableCaller:     true,
				DisableStacktrace: true,
				Sampling:          nil,
				Encoding:          "console",
				EncoderConfig: zapcore.EncoderConfig{
					TimeKey:          "T",
					LevelKey:         "L",
					NameKey:          "N",
					CallerKey:        "C",
					FunctionKey:      zapcore.OmitKey,
					MessageKey:       "M",
					StacktraceKey:    "S",
					LineEnding:       zapcore.DefaultLineEnding,
					EncodeLevel:      zapcore.CapitalLevelEncoder,
					EncodeTime:       zapcore.RFC3339TimeEncoder,
					EncodeDuration:   zapcore.StringDurationEncoder,
					EncodeCaller:     zapcore.ShortCallerEncoder,
					ConsoleSeparator: "\t|\t",
				},
				OutputPaths:      []string{"stderr"},
				ErrorOutputPaths: []string{"stderr"},
			}.Build()
			if err != nil {
				panic(err)
			}

			return logger.Sugar()
		}(),
		errorsContext:      errorsContext{},
		neutralizerContext: neutralizerContext{neutralizer: syntax_neutralizer.New(1)},
	}
}
