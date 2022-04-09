package context

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ FullContext = (*devContext)(nil)

type devContext struct {
	context.Context
	logger *zap.SugaredLogger
	errors []error
}

func (d *devContext) SetLogger(logger *zap.SugaredLogger) {
	d.logger = logger
}

func (d *devContext) Logger() *zap.SugaredLogger {
	return d.logger
}

func (d *devContext) AddError(err error) {
	d.errors = append(d.errors, err)
}

func (d *devContext) Errors() []error {
	return d.errors
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
		errors: nil,
	}
}
