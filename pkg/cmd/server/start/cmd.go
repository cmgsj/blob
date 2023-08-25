package start

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/cmgsj/blob/pkg/blob"
	cmdutil "github.com/cmgsj/blob/pkg/cmd/util"
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

type StartOptions struct {
	IOStreams cmdutil.IOStreams
	HttpAddr  string
	GrpcAddr  string
}

func NewStartOptions(streams cmdutil.IOStreams) *StartOptions {
	return &StartOptions{
		IOStreams: streams,
	}
}

func NewCmdStart(f cmdutil.Factory, streams cmdutil.IOStreams) *cobra.Command {
	o := NewStartOptions(streams)
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start blob server",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			stderr := o.IOStreams.Err
			cmdutil.CheckErr(o.Complete(f, cmd, args), stderr)
			cmdutil.CheckErr(o.Validate(), stderr)
			cmdutil.CheckErr(o.Run(f, cmd), stderr)
		},
	}
	return cmd
}

func (o *StartOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
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

func (o *StartOptions) Validate() error {
	return nil
}

func (o *StartOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
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
	ctx := context.Background()
	err := blobv1.RegisterBlobServiceHandlerClient(ctx, rmux, f.BlobServiceClient())
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", rmux)
	mux.Handle("/docs/", openapi.ServeDocs())

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
