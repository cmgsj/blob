package start

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	cmdutil "github.com/cmgsj/blob/pkg/cmd/util"
	"github.com/cmgsj/blob/pkg/interceptors"
	"github.com/spf13/cobra"

	"github.com/cmgsj/blob/pkg/blob"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/cmgsj/blob/pkg/openapi"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func NewCmdStart(streams cmdutil.IOStreams) *cobra.Command {
	o := NewStartOptions(streams)
	cmd := &cobra.Command{
		Use:  "start",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			stderr := o.IOStreams.Err
			cmdutil.CheckErr(o.Complete(cmd, args), stderr)
			cmdutil.CheckErr(o.Validate(), stderr)
			cmdutil.CheckErr(o.Run(), stderr)
		},
	}
	return cmd
}

func (o *StartOptions) Complete(cmd *cobra.Command, args []string) error {
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

func (o *StartOptions) Run() error {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.ServerUnaryLogger),
		grpc.StreamInterceptor(interceptors.ServerStreamLogger),
	)
	healthServer := health.NewServer()
	blobServer := blob.NewServer()

	reflection.Register(grpcServer)
	healthv1.RegisterHealthServer(grpcServer, healthServer)
	blobv1.RegisterBlobServiceServer(grpcServer, blobServer)

	healthServer.SetServingStatus(blob.ServiceName, healthv1.HealthCheckResponse_SERVING)

	rmux := runtime.NewServeMux()
	ctx := context.Background()
	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(interceptors.ClientUnaryLogger),
		grpc.WithStreamInterceptor(interceptors.ClientStreamLogger),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err := blobv1.RegisterBlobServiceHandlerFromEndpoint(ctx, rmux, o.GrpcAddr, opts)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", rmux)
	mux.Handle("/docs/", http.FileServer(http.FS(openapi.Docs())))

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
