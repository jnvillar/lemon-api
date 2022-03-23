package logger

import (
	"go.uber.org/zap"
)

type sugaredLogger struct {
	*zap.SugaredLogger
}

func NewSugaredLogger(zapSugaredLogger *zap.SugaredLogger) *sugaredLogger {
	return &sugaredLogger{
		zapSugaredLogger,
	}
}
