package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger

func init() {
	// Custom time encoder to format the timestamp
	timeFormat := func(timestamp time.Time, encode zapcore.PrimitiveArrayEncoder) {
		encode.AppendString(timestamp.Format(time.RFC3339))
	}

	// Custom encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeFormat,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
	}

	rawLogger, err := config.Build()
	if err != nil {
		panic(err)
	}
	log = rawLogger.Sugar()
	defer func() {
		if err := rawLogger.Sync(); err != nil {
			panic(err)
		}
	}()
}

// getLogger returns the global sugared logger instance
func getLogger() *zap.SugaredLogger {
	return log
}

// Infof logs an info message
func Infof(template string, args ...interface{}) {
	getLogger().Infof(template, args...)
}

// Warnf logs a warning message
func Warnf(template string, args ...interface{}) {
	getLogger().Warnf(template, args...)
}

// Errorf logs an error message
func Errorf(template string, args ...interface{}) {
	getLogger().Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	getLogger().Errorf(template, args...)
	os.Exit(1)
}
