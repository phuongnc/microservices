package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

func New() *Logger {
	lc := zap.NewProductionConfig()
	lc.EncoderConfig.MessageKey = "message"
	lc.EncoderConfig.TimeKey = "timestamp"
	lc.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	lc.EncoderConfig.CallerKey = ""
	lc.Encoding = "json"

	logger, err := lc.Build()
	if err != nil {
		panic(err)
	}

	return &Logger{logger: logger}
}

func (l *Logger) Info(message string) {
	l.logger.Info(message)
}

func (l *Logger) Error(message string, err error) {
	l.logger.Error(message, zap.Error(err))
}
