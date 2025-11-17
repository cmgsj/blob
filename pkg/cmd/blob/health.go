package blob

import (
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/worldline-go/grpc/health/v1/healthv1connect"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func NewCommandHealth() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Check blob service health",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			baseURL := viper.GetString("base-url")
			service := viper.GetString("service")

			healthClient := healthv1connect.NewHealthClient(http.DefaultClient, baseURL)

			request := &grpc_health_v1.HealthCheckRequest{
				Service: service,
			}

			response, err := healthClient.Check(ctx, connect.NewRequest(request))
			if err != nil {
				return err
			}

			if service != "" {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), service+": "+response.Msg.GetStatus().String())
			} else {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), response.Msg.GetStatus().String())
			}

			return nil
		},
	}

	cmd.Flags().String("base-url", "http://localhost:8080", "blob service base url")
	cmd.Flags().String("service", "", "service")

	_ = viper.BindPFlags(cmd.Flags())

	return cmd
}
