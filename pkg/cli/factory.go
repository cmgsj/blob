package cli

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/cmgsj/blob/pkg/interceptors"
)

type Factory interface {
	BlobServiceClient(ctx context.Context) (blobv1.BlobServiceClient, error)
	HealthClient(ctx context.Context) (healthv1.HealthClient, error)
}

type ConfigGetter interface {
	GetGRPCAddress() string
	GetHTTPAddress() string
}

func NewFactory(getter ConfigGetter) Factory {
	return &factory{
		ConfigGetter: getter,
	}
}

type factory struct {
	ConfigGetter
}

func (f *factory) BlobServiceClient(ctx context.Context) (blobv1.BlobServiceClient, error) {
	conn, err := f.dial(ctx)
	if err != nil {
		return nil, err
	}
	return blobv1.NewBlobServiceClient(conn), nil
}

func (f *factory) HealthClient(ctx context.Context) (healthv1.HealthClient, error) {
	conn, err := f.dial(ctx)
	if err != nil {
		return nil, err
	}
	return healthv1.NewHealthClient(conn), nil
}

func (f *factory) dial(ctx context.Context) (conn *grpc.ClientConn, err error) {
	logger := interceptors.NewLogger()
	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(logger.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(logger.StreamClientInterceptor()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	return grpc.DialContext(ctx, f.GetGRPCAddress(), opts...)
}
