package interceptors

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
)

func NewRecoverer() ServerInterceptor {
	opts := []recovery.Option{
		recovery.WithRecoveryHandlerContext(recoveryHandler),
	}
	return interceptor{
		us: recovery.UnaryServerInterceptor(opts...),
		ss: recovery.StreamServerInterceptor(opts...),
	}
}

func recoveryHandler(ctx context.Context, p any) (err error) {
	return fmt.Errorf("%v", p)
}
