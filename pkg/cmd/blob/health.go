package blob

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

func NewCommandHealth() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Health service client",
	}

	cmd.AddCommand(
		NewCommandHealthCheck(),
		NewCommandHealthList(),
		NewCommandHealthWatch(),
	)

	return cmd
}

func NewCommandHealthCheck() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check health status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			address := viper.GetString("address")
			service := viper.GetString("service")

			conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return err
			}

			healthClient := healthv1.NewHealthClient(conn)

			response, err := healthClient.Check(ctx, &healthv1.HealthCheckRequest{
				Service: service,
			})
			if err != nil {
				return err
			}

			if service != "" {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), service+":", response.GetStatus())
			} else {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), response.GetStatus())
			}

			return nil
		},
	}

	cmd.Flags().String("address", "localhost:2562", "server address")
	cmd.Flags().String("service", "", "service")

	_ = viper.BindPFlags(cmd.Flags())

	return cmd
}

func NewCommandHealthList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List health statuses",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			address := viper.GetString("address")

			conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return err
			}

			healthClient := healthv1.NewHealthClient(conn)

			response, err := healthClient.List(ctx, &healthv1.HealthListRequest{})
			if err != nil {
				return err
			}

			statuses := response.GetStatuses()

			tab := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 1, ' ', 0)

			for _, service := range slices.Sorted(maps.Keys(statuses)) {
				_, _ = fmt.Fprintln(tab, service+":\t", statuses[service].GetStatus())
			}

			err = tab.Flush()
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String("address", "localhost:2562", "server address")

	_ = viper.BindPFlags(cmd.Flags())

	return cmd
}

func NewCommandHealthWatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch health status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			address := viper.GetString("address")
			service := viper.GetString("service")

			conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return err
			}

			healthClient := healthv1.NewHealthClient(conn)

			stream, err := healthClient.Watch(ctx, &healthv1.HealthCheckRequest{
				Service: service,
			})
			if err != nil {
				return err
			}

			for {
				response, err := stream.Recv()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}

					return err
				}

				if service != "" {
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), service+":", response.GetStatus())
				} else {
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), response.GetStatus())
				}
			}

			return nil
		},
	}

	cmd.Flags().String("address", "localhost:2562", "server address")
	cmd.Flags().String("service", "", "service")

	_ = viper.BindPFlags(cmd.Flags())

	return cmd
}
