package util

import (
	"context"

	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/cmgsj/blob/pkg/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

type Factory interface {
	BlobServiceClient() blobv1.BlobServiceClient
	HealthClient() healthv1.HealthClient
}

type factory struct {
	blobServiceClient blobv1.BlobServiceClient
	healthClient      healthv1.HealthClient
}

func (f *factory) BlobServiceClient() blobv1.BlobServiceClient {
	return f.blobServiceClient
}

func (f *factory) HealthClient() healthv1.HealthClient {
	return f.healthClient
}

func NewFactory(ctx context.Context, address string) (Factory, error) {
	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(interceptors.ClientUnaryLogger),
		grpc.WithStreamInterceptor(interceptors.ClientStreamLogger),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.DialContext(ctx, address, opts...)
	if err != nil {
		return nil, err
	}
	return &factory{
		blobServiceClient: blobv1.NewBlobServiceClient(conn),
		healthClient:      healthv1.NewHealthClient(conn),
	}, nil
}
