package blob

import (
	"log/slog"
	"net"
	"net/http"

	"github.com/cmgsj/go-lib/swagger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/cmgsj/blob/pkg/blob"
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
			logger := interceptors.NewLogger()

			grpcServer := grpc.NewServer(
				grpc.UnaryInterceptor(logger.UnaryServerInterceptor()),
				grpc.StreamInterceptor(logger.StreamServerInterceptor()),
			)
			healthServer := health.NewServer()
			blobServer := blob.NewServer()

			reflection.Register(grpcServer)
			healthv1.RegisterHealthServer(grpcServer, healthServer)
			blobv1.RegisterBlobServiceServer(grpcServer, blobServer)

			healthServer.SetServingStatus(blob.ServiceName, healthv1.HealthCheckResponse_SERVING)

			rmux := runtime.NewServeMux()
			client, err := c.BlobServiceClient()
			if err != nil {
				return err
			}
			err = blobv1.RegisterBlobServiceHandlerClient(cmd.Context(), rmux, client)
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
				err = http.Serve(lis, mux)
				if err != nil {
					errch <- err
				}
			}()

			go func() {
				lis, err := net.Listen("tcp", c.GRPCAddress())
				if err != nil {
					errch <- err
				}
				slog.Info("started listening", "protocol", "grpc", "address", c.GRPCAddress())
				err = grpcServer.Serve(lis)
				if err != nil {
					errch <- err
				}
			}()

			for err := range errch {
				return err
			}

			return nil
		},
	}

	return cmd
}
