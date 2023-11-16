package start

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/cmgsj/blob/pkg/blob"
	"github.com/cmgsj/blob/pkg/cli"
	"github.com/cmgsj/blob/pkg/docs"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/cmgsj/blob/pkg/interceptors"
	"github.com/cmgsj/go-lib/openapi"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func NewCmdStart(f cli.Factory) *cobra.Command {
	o := NewStartOptions(f)
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start blob server",
		Args:  cobra.NoArgs,
		Run:   cli.Run(o),
	}
	return cmd
}

type StartOptions struct {
	cli.Factory
	HttpAddr string
	GrpcAddr string
}

func NewStartOptions(f cli.Factory) *StartOptions {
	return &StartOptions{
		Factory: f,
	}
}

func (o *StartOptions) Complete(ctx context.Context, cmd *cobra.Command, args []string) error {
	var err error
	o.HttpAddr, err = cmd.Flags().GetString("http-address")
	if err != nil {
		return err
	}
	o.GrpcAddr, err = cmd.Flags().GetString("grpc-address")
	if err != nil {
		return err
	}
	return nil
}

func (o *StartOptions) Validate(ctx context.Context) error {
	return nil
}

func (o *StartOptions) Run(ctx context.Context) error {
	logger := interceptors.NewLogger()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logger.UnaryServerInterceptor()),
		grpc.StreamInterceptor(logger.StreamServerInterceptor()),
	)
	healthServer := health.NewServer()
	blobServer := blob.NewServer()

	reflection.RegisterV1(grpcServer)
	healthv1.RegisterHealthServer(grpcServer, healthServer)
	blobv1.RegisterBlobServiceServer(grpcServer, blobServer)

	healthServer.SetServingStatus(blob.ServiceName, healthv1.HealthCheckResponse_SERVING)

	rmux := runtime.NewServeMux()
	client, err := o.BlobServiceClient(ctx)
	if err != nil {
		return err
	}
	err = blobv1.RegisterBlobServiceHandlerClient(ctx, rmux, client)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", rmux)
	mux.Handle("/docs/", openapi.ServeDocs("/docs/", docs.OpenapiSchema))

	errch := make(chan error)

	go func() {
		lis, err := net.Listen("tcp", o.HttpAddr)
		if err != nil {
			errch <- err
		}
		slog.Info("started listening", "protocol", "http", "address", o.HttpAddr)
		err = http.Serve(lis, mux)
		if err != nil {
			errch <- err
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", o.GrpcAddr)
		if err != nil {
			errch <- err
		}
		slog.Info("started listening", "protocol", "grpc", "address", o.GrpcAddr)
		err = grpcServer.Serve(lis)
		if err != nil {
			errch <- err
		}
	}()

	for err := range errch {
		return err
	}

	return nil
}
