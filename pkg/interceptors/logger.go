package interceptors

import (
	"context"
	"log/slog"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func NewLogger() FullInterceptor {
	log := SlogAdapter(slog.Default())
	opts := []logging.Option{
		logging.WithCodes(logging.DefaultErrorToCode),
		logging.WithDurationField(logging.DefaultDurationToFields),
		logging.WithLevels(logging.DefaultClientCodeToLevel),
		logging.WithTimestampFormat(time.RFC3339),
		logging.WithFieldsFromContext(logging.ExtractFields),
	}
	return interceptor{
		us: logging.UnaryServerInterceptor(log, opts...),
		ss: logging.StreamServerInterceptor(log, opts...),
		uc: logging.UnaryClientInterceptor(log, opts...),
		sc: logging.StreamClientInterceptor(log, opts...),
	}
}

func SlogAdapter(log *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		log.Log(ctx, slog.Level(level), msg, fields...)
	})
}
