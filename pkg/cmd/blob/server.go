package blob

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"runtime/debug"
	"syscall"

	validate "buf.build/go/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/riza-io/grpc-go/credentials/basic"
	"github.com/riza-io/grpc-go/credentials/bearer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	reflectionv1 "google.golang.org/grpc/reflection/grpc_reflection_v1"
	reflectionv1alpha "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"

	"github.com/cmgsj/blob/pkg/blob"
	"github.com/cmgsj/blob/pkg/blob/storage"
	"github.com/cmgsj/blob/pkg/blob/storage/driver"
	"github.com/cmgsj/blob/pkg/blob/storage/driver/azblob"
	"github.com/cmgsj/blob/pkg/blob/storage/driver/gcs"
	"github.com/cmgsj/blob/pkg/blob/storage/driver/memory"
	"github.com/cmgsj/blob/pkg/blob/storage/driver/minio"
	"github.com/cmgsj/blob/pkg/blob/storage/driver/s3"
	blobv1 "github.com/cmgsj/blob/pkg/proto/blob/v1"
)

func NewCommandServer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start blob server",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
			defer cancel()

			address := viper.GetString("address")

			blobStorageDriver, err := newBlobStorageDriver(ctx)
			if err != nil {
				return err
			}

			blobStorage, err := storage.NewStorage(ctx, blobStorageDriver)
			if err != nil {
				return err
			}

			blobServer := blob.NewServer(blobStorage)

			healthServer := health.NewServer()

			healthServer.SetServingStatus(blobv1.BlobService_ServiceDesc.ServiceName, healthv1.HealthCheckResponse_SERVING)
			healthServer.SetServingStatus(healthv1.Health_ServiceDesc.ServiceName, healthv1.HealthCheckResponse_SERVING)
			healthServer.SetServingStatus(reflectionv1.ServerReflection_ServiceDesc.ServiceName, healthv1.HealthCheckResponse_SERVING)
			healthServer.SetServingStatus(reflectionv1alpha.ServerReflection_ServiceDesc.ServiceName, healthv1.HealthCheckResponse_SERVING)

			logger, err := newLogger()
			if err != nil {
				return err
			}

			credentials, err := newServerTransportCredentials()
			if err != nil {
				return err
			}

			recoveryHandlerFunc := newRecoveryHandlerFunc(logger)
			authFunc := newServerAuthFunc()
			authMatcher := newAuthMatcher()
			validator := newValidator()

			grpcServer := grpc.NewServer(
				grpc.Creds(credentials),
				grpc.ChainUnaryInterceptor(
					recovery.UnaryServerInterceptor(recovery.WithRecoveryHandlerContext(recoveryHandlerFunc)),
					logging.UnaryServerInterceptor(logger),
					selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFunc), authMatcher),
					protovalidate.UnaryServerInterceptor(validator),
				),
				grpc.ChainStreamInterceptor(
					recovery.StreamServerInterceptor(recovery.WithRecoveryHandlerContext(recoveryHandlerFunc)),
					logging.StreamServerInterceptor(logger),
					selector.StreamServerInterceptor(auth.StreamServerInterceptor(authFunc), authMatcher),
					protovalidate.StreamServerInterceptor(validator),
				),
			)

			blobv1.RegisterBlobServiceServer(grpcServer, blobServer)
			healthv1.RegisterHealthServer(grpcServer, healthServer)
			reflection.Register(grpcServer)

			listenConfig := &net.ListenConfig{}

			listener, err := listenConfig.Listen(ctx, "tcp", address)
			if err != nil {
				return err
			}

			go func() {
				logger.Log(ctx, logging.LevelInfo, "starting grpc server", "address", listener.Addr())

				err = grpcServer.Serve(listener)
				if err != nil {
					logger.Log(ctx, logging.LevelError, "failed to run grpc server", "error", err)
				}
			}()

			<-ctx.Done()

			grpcServer.GracefulStop()
			grpcServer.Stop()

			return nil
		},
	}

	cmd.Flags().String("driver-type", "memory", "driver type")
	cmd.Flags().String("driver-uri", "", "driver uri")

	return cmd
}

func newBlobStorageDriver(ctx context.Context) (driver.Driver, error) {
	driverType := viper.GetString("driver-type")
	driverURI := viper.GetString("driver-uri")

	switch driverType {
	case "memory":
		return memory.NewDriver(), nil

	case "azblob":
		return azblob.NewDriver(ctx, azblob.DriverOptions{
			URI: driverURI,
		})

	case "gcs":
		return gcs.NewDriver(ctx, gcs.DriverOptions{
			URI: driverURI,
		})

	case "s3":
		return s3.NewDriver(ctx, s3.DriverOptions{
			URI: driverURI,
		})

	case "minio":
		return minio.NewDriver(ctx, minio.DriverOptions{
			URI: driverURI,
		})

	default:
		return nil, fmt.Errorf("unknown blob storage driver %q", driverType)
	}
}

func newRecoveryHandlerFunc(logger logging.Logger) recovery.RecoveryHandlerFuncContext {
	return func(ctx context.Context, p any) error {
		logger.Log(ctx, logging.LevelError, "recovered from panic", "panic", p, "stack", debug.Stack())

		return status.Errorf(codes.Internal, "%s", p)
	}
}

func newBearerAuthFunc(authToken string) auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		token, err := bearer.TokenFromContext(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid bearer auth header")
		}

		if string(token) != authToken {
			return nil, status.Error(codes.Unauthenticated, "invalid bearer auth token")
		}

		return ctx, nil
	}
}

func newBasicAuthFunc(authUsername, authPassword string) auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		credentials, err := basic.CredentialsFromContext(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid basic auth header")
		}

		if credentials.UserID != authUsername || credentials.Password != authPassword {
			return nil, status.Error(codes.Unauthenticated, "invalid basic auth credentials")
		}

		return ctx, nil
	}
}

func newInsecureAuthFunc() auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		return ctx, nil
	}
}

func newAuthMatcher() selector.Matcher {
	return selector.MatchFunc(func(ctx context.Context, callMeta interceptors.CallMeta) bool {
		switch callMeta.Service {
		case
			healthv1.Health_ServiceDesc.ServiceName,
			reflectionv1.ServerReflection_ServiceDesc.ServiceName,
			reflectionv1alpha.ServerReflection_ServiceDesc.ServiceName:
			return false

		default:
			return true
		}
	})
}

func newValidator() validate.Validator {
	return validate.GlobalValidator
}
