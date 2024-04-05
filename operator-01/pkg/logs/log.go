package logs

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// set logging
// todo: demonstrate modified logging from main.go
func LogConfig() *zap.Logger {
	config := zap.NewProductionConfig()
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info" // Default LogLevel
	}
	var zapLogLevel zapcore.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		zapLogLevel = zapcore.DebugLevel
	case "info":
		zapLogLevel = zapcore.InfoLevel
	case "warn":
		zapLogLevel = zapcore.WarnLevel
	case "error":
		zapLogLevel = zapcore.ErrorLevel
	default:
		zapLogLevel = zapcore.InfoLevel
	}
	config.EncoderConfig.NameKey = "logger"
	config.EncoderConfig.MessageKey = "log"
	config.EncoderConfig.StacktraceKey = "stacktrace"
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.LevelKey = "level"
	config.Level.SetLevel(zapLogLevel)
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	return logger
}

func Logger() *zap.SugaredLogger {
	logger := LogConfig()
	sugar := logger.Sugar()
	return sugar
}

func LoggerCtrl() *zap.Logger {
	return LogConfig()
}
