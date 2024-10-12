package blob

import (
	"fmt"

	"github.com/spf13/cobra"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/cmgsj/blob/pkg/blob"
	"github.com/cmgsj/blob/pkg/cli"
)

func NewCmdServerHealth(c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "health-check blob server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			req := &healthv1.HealthCheckRequest{
				Service: blob.ServiceName,
			}

			client, err := c.HealthClient()
			if err != nil {
				return err
			}

			resp, err := client.Check(cmd.Context(), req)
			if err != nil {
				return err
			}

			service := req.GetService()

			status := healthv1.HealthCheckResponse_ServingStatus_name[int32(resp.GetStatus())]

			fmt.Printf("%s: %s\n", service, status)

			return nil
		},
	}

	return cmd
}
