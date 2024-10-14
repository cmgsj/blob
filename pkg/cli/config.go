package cli

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/cmgsj/blob/pkg/interceptors"
)

type ConfigOptions struct {
	GRPCAddress string
	HTTPAddress string
}

type Config struct {
	opts ConfigOptions
}

func NewConfig(opts ConfigOptions) *Config {
	return &Config{
		opts: opts,
	}
}

func (c *Config) GRPCAddress() string {
	return c.opts.GRPCAddress
}

func (c *Config) HTTPAddress() string {
	return c.opts.HTTPAddress
}

func (c *Config) BlobServiceClient() (blobv1.BlobServiceClient, error) {
	conn, err := c.grpcClientDial()
	if err != nil {
		return nil, err
	}

	return blobv1.NewBlobServiceClient(conn), nil
}

func (c *Config) HealthClient() (healthv1.HealthClient, error) {
	conn, err := c.grpcClientDial()
	if err != nil {
		return nil, err
	}

	return healthv1.NewHealthClient(conn), nil
}

func (c *Config) grpcClientDial() (conn *grpc.ClientConn, err error) {
	logger := interceptors.NewLogger()

	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(logger.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(logger.StreamClientInterceptor()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	return grpc.NewClient(c.opts.GRPCAddress, opts...)
}
