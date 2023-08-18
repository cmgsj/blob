package util

import (
	"context"

	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/cmgsj/blob/pkg/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Factory interface {
	BlobServiceClient() blobv1.BlobServiceClient
}

type factory struct {
	blobServiceClient blobv1.BlobServiceClient
}

func (f *factory) BlobServiceClient() blobv1.BlobServiceClient {
	return f.blobServiceClient
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
	}, nil
}
