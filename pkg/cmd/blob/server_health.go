package blob

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/cmgsj/blob/pkg/cli"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

func NewCmdServerHealth(c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "health-check blob server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			service := viper.GetString("service")

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

			status := healthv1.HealthCheckResponse_ServingStatus_name[int32(resp.GetStatus())]

			if service != "" {
				fmt.Println(status + " (" + service + ")")
			} else {
				fmt.Println(status)
			}

			return nil
		},
	}

	cmd.Flags().String("service", blobv1.BlobService_ServiceDesc.ServiceName, "grpc service")

	viper.BindPFlags(cmd.Flags())

	return cmd
}
