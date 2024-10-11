package interceptors

import (
	"context"
	"log/slog"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func NewLogger() FullInterceptor {
	logger := SlogAdapter(slog.Default())

	opts := []logging.Option{
		logging.WithCodes(logging.DefaultErrorToCode),
		logging.WithDurationField(logging.DefaultDurationToFields),
		logging.WithLevels(logging.DefaultClientCodeToLevel),
		logging.WithTimestampFormat(time.RFC3339),
		logging.WithFieldsFromContext(logging.ExtractFields),
	}

	return FullInterceptor{
		clientInterceptor: clientInterceptor{
			unary:  logging.UnaryClientInterceptor(logger, opts...),
			stream: logging.StreamClientInterceptor(logger, opts...),
		},
		serverInterceptor: serverInterceptor{
			unary:  logging.UnaryServerInterceptor(logger, opts...),
			stream: logging.StreamServerInterceptor(logger, opts...),
		},
	}
}

func SlogAdapter(log *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		log.Log(ctx, slog.Level(level), msg, fields...)
	})
}
