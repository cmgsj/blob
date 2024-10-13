package blob

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/cmgsj/go-lib/swagger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	reflectionv1 "google.golang.org/grpc/reflection/grpc_reflection_v1"
	reflectionv1alpha "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	blobserver "github.com/cmgsj/blob/pkg/blob/server"
	"github.com/cmgsj/blob/pkg/blob/storage"
	"github.com/cmgsj/blob/pkg/blob/storage/gcs"
	"github.com/cmgsj/blob/pkg/blob/storage/memory"
	"github.com/cmgsj/blob/pkg/blob/storage/minio"
	"github.com/cmgsj/blob/pkg/blob/storage/mongodb"
	"github.com/cmgsj/blob/pkg/cli"
	"github.com/cmgsj/blob/pkg/docs"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/cmgsj/blob/pkg/interceptors"
)

func NewCmdServerStart(c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start blob server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			blobStorage, err := newBlobStorage(ctx)
			if err != nil {
				return err
			}

			blobServer := blobserver.NewServer(blobStorage)

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

	cmd.Flags().String("storage", "memory", "blob storage type")

	cmd.Flags().String("gcs-uri", "", "gcs uri")

	cmd.Flags().String("minio-address", "", "minio address")
	cmd.Flags().String("minio-access-key", "", "minio access key")
	cmd.Flags().String("minio-secret-key", "", "minio secret key")
	cmd.Flags().Bool("minio-secure", false, "minio secure")
	cmd.Flags().String("minio-bucket", "", "minio bucket")
	cmd.Flags().String("minio-object-prefix", "", "minio object prefix")

	cmd.Flags().String("mongodb-uri", "", "mongodb uri")

	viper.BindPFlags(cmd.Flags())

	return cmd
}

func newBlobStorage(ctx context.Context) (storage.Storage, error) {
	storage := viper.GetString("storage")

	switch storage {
	case "memory":
		return memory.NewStorage(), nil

	case "gcs":
		return gcs.NewStorage(ctx, gcs.StorageOptions{
			URI: viper.GetString("gcs-uri"),
		})

	case "mongodb":
		return mongodb.NewStorage(ctx, mongodb.StorageOptions{
			URI: viper.GetString("mongodb-uri"),
		})

	case "minio":
		return minio.NewStorage(ctx, minio.StorageOptions{
			Address:      viper.GetString("minio-address"),
			AccessKey:    viper.GetString("minio-access-key"),
			SecretKey:    viper.GetString("minio-secret-key"),
			Secure:       viper.GetBool("minio-secure"),
			Bucket:       viper.GetString("minio-bucket"),
			ObjectPrefix: viper.GetString("minio-object-prefix"),
		})
	}

	return nil, fmt.Errorf("unknown blob storage %q", storage)
}
