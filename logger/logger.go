package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Sugar *sugaredLogger

func init() {
	config, _ := zap.NewProductionConfig().Build()
	Sugar = NewSugaredLogger(config.Sugar())
}

func Initialize(level, version string) error {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(getZapLevel(level))
	logger, err := config.Build()
	if err != nil {
		return err
	}

	Sugar = NewSugaredLogger(logger.Sugar())
	return nil
}

func getZapLevel(level string) zapcore.Level {
	levelMap := map[string]zapcore.Level{
		"debug":   zapcore.DebugLevel,
		"info":    zapcore.InfoLevel,
		"error":   zapcore.ErrorLevel,
		"warning": zapcore.WarnLevel,
		"fatal":   zapcore.FatalLevel,
		"panic":   zapcore.PanicLevel,
	}
	return levelMap[level]
}
