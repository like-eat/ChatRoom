package log

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func GetZapLogger() *zap.Logger {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	return logger
}
