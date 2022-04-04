package logger

import (
	"os"

	"go.uber.org/zap"
)

func New() *zap.SugaredLogger {
	if os.Getenv("DEBUG") == "1" {
		logger, _ := zap.NewDevelopment()
		return logger.Sugar()
	}

	logger, _ := zap.NewProduction()
	return logger.Sugar()
}
