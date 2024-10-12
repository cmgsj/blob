package interceptors

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
)

func NewRecovery() ServerInterceptor {
	opts := []recovery.Option{
		recovery.WithRecoveryHandlerContext(recoveryHandler),
	}

	return ServerInterceptor{
		serverInterceptor: serverInterceptor{
			unary:  recovery.UnaryServerInterceptor(opts...),
			stream: recovery.StreamServerInterceptor(opts...),
		},
	}
}

func recoveryHandler(ctx context.Context, p any) error {
	return fmt.Errorf("%v", p)
}
