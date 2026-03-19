package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *Logger

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger() *Logger {
	if globalLogger != nil {
		return globalLogger
	}

	config := zap.NewDevelopmentConfig()
	config.DisableStacktrace = true
	config.EncoderConfig.FunctionKey = "F"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	logger, _ := config.Build()
	defer logger.Sync()
	sugar := logger.Sugar()
	globalLogger = &Logger{sugar}
	return globalLogger
}

func GetLogger() *Logger {
	return NewLogger()
}
