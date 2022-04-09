package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() *zap.SugaredLogger {
	if os.Getenv("DEBUG") == "1" {
		logger, _ := zap.Config{
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

		return logger.Sugar()
	}

	return zap.NewNop().Sugar()
}
