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

	// Get log level from environment variable
	logLevel := getLogLevelFromEnv()

	config := zap.Config{
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		Level:            zap.NewAtomicLevelAt(logLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	rawLogger, err := config.Build()
	if err != nil {
		panic(err)
	}
	log = rawLogger.Sugar()
}

// getLogLevelFromEnv reads the log level from the environment variable LOG_LEVEL
func getLogLevelFromEnv() zapcore.Level {
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
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
	getLogger().Fatalf(template, args...)
	os.Exit(1)
}
