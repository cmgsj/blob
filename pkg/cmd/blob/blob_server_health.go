package blob

import (
	"fmt"

	"github.com/spf13/cobra"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/cmgsj/blob/pkg/cli"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

func NewCmdServerHealth(c *cli.Config) *cobra.Command {
	var service string = blobv1.BlobService_ServiceDesc.ServiceName

	cmd := &cobra.Command{
		Use:   "health",
		Short: "health-check blob server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			healthClient, err := c.HealthClient()
			if err != nil {
				return err
			}

			resp, err := healthClient.Check(ctx, &healthv1.HealthCheckRequest{
				Service: service,
			})
			if err != nil {
				return err
			}

			fmt.Printf("%s: %s\n", service, healthv1.HealthCheckResponse_ServingStatus_name[int32(resp.GetStatus())])

			return nil
		},
	}

	cmd.Flags().StringVar(&service, "service", service, "grpc service")

	return cmd
}
