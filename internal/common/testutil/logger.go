package testutil

import (
	"bytes"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLoggerWithBuffer(buf *bytes.Buffer) *zap.Logger {
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	writer := zapcore.Lock(zapcore.AddSync(buf))
	enabler := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.ErrorLevel
	})
	return zap.New(zapcore.NewCore(encoder, writer, enabler))
}
