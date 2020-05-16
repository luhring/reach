package reachlog

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the interface that wraps methods for Reach to use during logging
type Logger interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
	Fatal(message string, args ...interface{})
}

// Level indicates a specified logging level
type Level string

// The levels of logging that the Logger supports
const (
	LevelNone  Level = "none"
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelFatal Level = "fatal"
)

type zapWrapper struct {
	sugar *zap.SugaredLogger
}

// New returns a new, implemented instance of a Logger type
func New(minLevel Level) Logger {
	var l zapcore.Level
	outputPaths := []string{"stderr"}

	switch minLevel {
	case LevelNone:
		l = zapcore.FatalLevel
		outputPaths = nil
	case LevelDebug:
		l = zapcore.DebugLevel
	case LevelInfo:
		l = zapcore.InfoLevel
	case LevelWarn:
		l = zapcore.WarnLevel
	case LevelError:
		l = zapcore.ErrorLevel
	case LevelFatal:
		l = zapcore.FatalLevel
	default:
		l = zapcore.FatalLevel
		outputPaths = nil
	}
	var cfg zap.Config
	cfg.Level = zap.NewAtomicLevel()
	cfg.Level.SetLevel(l)
	cfg.Encoding = "console"
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.OutputPaths = outputPaths

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	sugar := logger.Sugar()
	return &zapWrapper{sugar: sugar}
}

// Debug logs the received message and parameters at the DEBUG logging level
func (w zapWrapper) Debug(message string, args ...interface{}) {
	w.sugar.Debugw(message, args...)
}

// Info logs the received message and parameters at the INFO logging level
func (w zapWrapper) Info(message string, args ...interface{}) {
	w.sugar.Infow(message, args...)
}

// Warn logs the received message and parameters at the WARN logging level
func (w zapWrapper) Warn(message string, args ...interface{}) {
	w.sugar.Warnw(message, args...)
}

// Error logs the received message and parameters at the ERROR logging level
func (w zapWrapper) Error(message string, args ...interface{}) {
	w.sugar.Errorw(message, args...)
}

// Fatal logs the received message and parameters at the FATAL logging level
func (w zapWrapper) Fatal(message string, args ...interface{}) {
	w.sugar.Fatalw(message, args...)
}
