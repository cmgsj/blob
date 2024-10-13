package blob

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"

	"github.com/cmgsj/go-lib/swagger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	reflectionv1 "google.golang.org/grpc/reflection/grpc_reflection_v1"
	reflectionv1alpha "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"github.com/cmgsj/blob/pkg/blob"
	"github.com/cmgsj/blob/pkg/blob/storage"
	"github.com/cmgsj/blob/pkg/cli"
	"github.com/cmgsj/blob/pkg/docs"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/cmgsj/blob/pkg/interceptors"
)

func NewCmdServerStart(c *cli.Config) *cobra.Command {
	var storage string = ":memory:"

	cmd := &cobra.Command{
		Use:   "start",
		Short: "start blob server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			blobStorage, err := parseBlobStorage(ctx, storage)
			if err != nil {
				return err
			}

			blobServer := blob.NewServer(blobStorage)

			healthServer := health.NewServer()
			healthServer.SetServingStatus(blobv1.BlobService_ServiceDesc.ServiceName, healthv1.HealthCheckResponse_SERVING)
			healthServer.SetServingStatus(healthv1.Health_ServiceDesc.ServiceName, healthv1.HealthCheckResponse_SERVING)
			healthServer.SetServingStatus(reflectionv1.ServerReflection_ServiceDesc.ServiceName, healthv1.HealthCheckResponse_SERVING)
			healthServer.SetServingStatus(reflectionv1alpha.ServerReflection_ServiceDesc.ServiceName, healthv1.HealthCheckResponse_SERVING)

			logger := interceptors.NewLogger()

			grpcServer := grpc.NewServer(
				grpc.UnaryInterceptor(logger.UnaryServerInterceptor()),
				grpc.StreamInterceptor(logger.StreamServerInterceptor()),
			)

			blobv1.RegisterBlobServiceServer(grpcServer, blobServer)
			healthv1.RegisterHealthServer(grpcServer, healthServer)
			reflection.Register(grpcServer)

			rmux := runtime.NewServeMux()

			blobClient, err := c.BlobServiceClient()
			if err != nil {
				return err
			}

			err = blobv1.RegisterBlobServiceHandlerClient(ctx, rmux, blobClient)
			if err != nil {
				return err
			}

			mux := http.NewServeMux()

			mux.Handle("/", rmux)
			mux.Handle("/docs/", swagger.Docs("/docs/", docs.SwaggerSchema()))

			errch := make(chan error)

			go func() {
				lis, err := net.Listen("tcp", c.HTTPAddress())
				if err != nil {
					errch <- err
				}
				slog.Info("started listening", "protocol", "http", "address", c.HTTPAddress())
				errch <- http.Serve(lis, mux)
			}()

			go func() {
				lis, err := net.Listen("tcp", c.GRPCAddress())
				if err != nil {
					errch <- err
				}
				slog.Info("started listening", "protocol", "grpc", "address", c.GRPCAddress())
				errch <- grpcServer.Serve(lis)
			}()

			return <-errch
		},
	}

	cmd.Flags().StringVar(&storage, "storage", storage, "blob storage url")

	return cmd
}

func parseBlobStorage(ctx context.Context, uri string) (blob.Storage, error) {
	if uri == ":memory:" {
		return storage.NewInMemoryStorage(), nil
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "gs":
		return storage.NewGoogleCloudStorage(ctx, uri)

	case "mongodb", "mongodb+srv":
		return storage.NewMongoDBStorage(ctx, uri)
	}

	return nil, fmt.Errorf("unknown storage %q", uri)
}
