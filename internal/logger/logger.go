package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	config := zap.NewProductionConfig()

	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.TimeKey = "timestamp"
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	config.EncoderConfig = encodeConfig

	logger, _ := config.Build()

	return logger
}

type zapKey struct{}

func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, zapKey{}, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	return ctx.Value(zapKey{}).(*zap.Logger)
}
