package log

import (
	"go.uber.org/zap"
)

func NewZapLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	return logger
}

func New() Logger {
	return &zapLogger{
		logger: NewZapLogger(),
	}
}

type zapLogger struct {
	logger *zap.Logger
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, zapFields(fields)...)
}

func zapFields(in []Field) []zap.Field {
	out := make([]zap.Field, 0, len(in))
	for _, f := range in {
		switch f.kind {
		case kindString:
			out = append(out, zap.String(f.key, f.value.(string)))
		case kindError:
			out = append(out, zap.Error(f.value.(error)))
		}

	}
	return out
}
