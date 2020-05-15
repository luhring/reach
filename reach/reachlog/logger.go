package reachlog

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
	Fatal(message string, args ...interface{})
}

type zapWrapper struct {
	sugar *zap.SugaredLogger
}

func New() Logger {
	var cfg zap.Config
	cfg.Level = zap.NewAtomicLevel()
	cfg.Level.SetLevel(zapcore.InfoLevel)
	cfg.Encoding = "console"
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.OutputPaths = []string{"stdout"}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	sugar := logger.Sugar()
	return &zapWrapper{sugar: sugar}
}

func (w zapWrapper) Debug(message string, args ...interface{}) {
	w.sugar.Debugw(message, args...)
}

func (w zapWrapper) Info(message string, args ...interface{}) {
	w.sugar.Infow(message, args...)
}

func (w zapWrapper) Warn(message string, args ...interface{}) {
	w.sugar.Warnw(message, args...)
}

func (w zapWrapper) Error(message string, args ...interface{}) {
	w.sugar.Errorw(message, args...)
}

func (w zapWrapper) Fatal(message string, args ...interface{}) {
	w.sugar.Fatalw(message, args...)
}
