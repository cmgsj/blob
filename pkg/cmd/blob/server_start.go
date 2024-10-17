package blob

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"net"
	"net/http"
	"slices"
	"strings"

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
	blobstorage "github.com/cmgsj/blob/pkg/blob/storage"
	blobstorageriver "github.com/cmgsj/blob/pkg/blob/storage/driver"
	"github.com/cmgsj/blob/pkg/blob/storage/driver/azblob"
	"github.com/cmgsj/blob/pkg/blob/storage/driver/gcs"
	"github.com/cmgsj/blob/pkg/blob/storage/driver/memory"
	"github.com/cmgsj/blob/pkg/blob/storage/driver/minio"
	"github.com/cmgsj/blob/pkg/blob/storage/driver/s3"
	"github.com/cmgsj/blob/pkg/cli"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/cmgsj/blob/pkg/interceptors"
	"github.com/cmgsj/blob/swagger"
)

func NewCmdServerStart(c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start blob server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			blobStorageDriver, err := newBlobStorageDriver(ctx)
			if err != nil {
				return err
			}

			blobStorage, err := blobstorage.NewStorage(ctx, blobStorageDriver)
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
			mux.Handle("/docs/", http.StripPrefix("/docs", swagger.Handler()))

			httpServer := &http.Server{
				Handler: mux,
			}

			slog.Info("starting blob server", "grpc_address", c.GRPCAddress(), "http_address", c.HTTPAddress())

			errch := make(chan error)

			go func() {
				lis, err := net.Listen("tcp", c.GRPCAddress())
				if err != nil {
					errch <- err
				}
				errch <- grpcServer.Serve(lis)
			}()

			go func() {
				lis, err := net.Listen("tcp", c.HTTPAddress())
				if err != nil {
					errch <- err
				}
				errch <- httpServer.Serve(lis)
			}()

			return <-errch
		},
	}

	cmd.Flags().String("azblob-uri", "", "azblob uri")
	cmd.Flags().String("azblob-account-name", "", "azblob account name")
	cmd.Flags().String("azblob-account-key", "", "azblob account key")

	cmd.Flags().String("gcs-uri", "", "gcs uri")

	cmd.Flags().String("s3-uri", "", "s3 uri")
	cmd.Flags().String("s3-access-key", "", "s3 access key")
	cmd.Flags().String("s3-secret-key", "", "s3 secret key")

	cmd.Flags().String("minio-address", "", "minio address")
	cmd.Flags().String("minio-access-key", "", "minio access key")
	cmd.Flags().String("minio-secret-key", "", "minio secret key")
	cmd.Flags().String("minio-bucket", "", "minio bucket")
	cmd.Flags().String("minio-object-prefix", "", "minio object prefix")
	cmd.Flags().Bool("minio-secure", false, "minio secure")

	viper.BindPFlags(cmd.Flags())

	return cmd
}

func newBlobStorageDriver(ctx context.Context) (blobstorageriver.Driver, error) {
	driverTypes := []string{
		"azblob",
		"gcs",
		"s3",
		"minio",
	}

	driverTypeSet := make(map[string]struct{})

	for _, key := range viper.AllKeys() {
		for _, driverType := range driverTypes {
			if strings.HasPrefix(key, driverType) && viper.IsSet(key) {
				driverTypeSet[driverType] = struct{}{}
				break
			}
		}
	}

	driverTypes = slices.Collect(maps.Keys(driverTypeSet))

	if len(driverTypes) > 1 {
		return nil, fmt.Errorf("multiple blob storage drivers set: [%s]", strings.Join(driverTypes, ", "))
	}

	if len(driverTypes) == 0 {
		driverTypes = []string{"memory"}
	}

	driverType := driverTypes[0]

	var driver blobstorageriver.Driver
	var err error

	switch driverType {
	case "memory":
		driver = memory.NewDriver()

	case "azblob":
		driver, err = azblob.NewDriver(ctx, azblob.DriverOptions{
			URI:         viper.GetString("azblob-uri"),
			AccountName: viper.GetString("azblob-account-name"),
			AccountKey:  viper.GetString("azblob-account-key"),
		})

	case "gcs":
		driver, err = gcs.NewDriver(ctx, gcs.DriverOptions{
			URI: viper.GetString("gcs-uri"),
		})

	case "s3":
		driver, err = s3.NewDriver(ctx, s3.DriverOptions{
			URI:       viper.GetString("s3-uri"),
			AccessKey: viper.GetString("s3-access-key"),
			SecretKey: viper.GetString("s3-secret-key"),
		})

	case "minio":
		driver, err = minio.NewDriver(ctx, minio.DriverOptions{
			Address:      viper.GetString("minio-address"),
			AccessKey:    viper.GetString("minio-access-key"),
			SecretKey:    viper.GetString("minio-secret-key"),
			Bucket:       viper.GetString("minio-bucket"),
			ObjectPrefix: viper.GetString("minio-object-prefix"),
			Secure:       viper.GetBool("minio-secure"),
		})

	default:
		return nil, fmt.Errorf("unknown blob storage driver %q", driverType)
	}
	if err != nil {
		return nil, err
	}

	slog.Info("using blob storage driver", "type", driverType)

	return driver, nil
}
