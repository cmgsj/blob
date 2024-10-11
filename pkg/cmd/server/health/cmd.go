package health

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/cmgsj/blob/pkg/blob"
	"github.com/cmgsj/blob/pkg/cli"
)

func NewCmdHealth(f cli.Factory) *cobra.Command {
	o := NewWriteOptions(f)
	cmd := &cobra.Command{
		Use:   "health",
		Short: "health-check blob server",
		Args:  cobra.NoArgs,
		Run:   cli.Run(o),
	}
	return cmd
}

type HealthOptions struct {
	cli.Factory
	Request *healthv1.HealthCheckRequest
}

func NewWriteOptions(f cli.Factory) *HealthOptions {
	return &HealthOptions{
		Factory: f,
		Request: &healthv1.HealthCheckRequest{
			Service: blob.ServiceName,
		},
	}
}

func (o *HealthOptions) Complete(ctx context.Context, cmd *cobra.Command, args []string) error {
	return nil
}

func (o *HealthOptions) Validate(ctx context.Context) error {
	return nil
}

func (o *HealthOptions) Run(ctx context.Context) error {
	client, err := o.HealthClient(ctx)
	if err != nil {
		return err
	}
	resp, err := client.Check(ctx, o.Request)
	if err != nil {
		return err
	}
	service := o.Request.GetService()
	status := healthv1.HealthCheckResponse_ServingStatus_name[int32(resp.GetStatus())]
	fmt.Printf("%s: %s\n", service, status)
	return nil
}
