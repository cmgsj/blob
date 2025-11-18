package blob

import (
	"cmp"
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
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

			_, _ = fmt.Fprintln(cmd.OutOrStdout(), response.GetStatus().String())

			return nil
		},
	}

	cmd.Flags().String("service", "", "service")

	return cmd
}

func NewCommandHealthList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List health statuses",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
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

			tab := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 1, ' ', 0)

			for _, service := range slices.Sorted(maps.Keys(response.GetStatuses())) {
				_, _ = fmt.Fprintln(tab, cmp.Or(service, "*")+":\t"+response.GetStatuses()[service].GetStatus().String())
			}

			err = tab.Flush()
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func NewCommandHealthWatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch health status",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
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

				_, _ = fmt.Fprintln(cmd.OutOrStdout(), response.GetStatus().String())
			}

			return nil
		},
	}

	cmd.Flags().String("service", "", "service")

	return cmd
}
