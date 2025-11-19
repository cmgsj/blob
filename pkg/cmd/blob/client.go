package blob

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	blobv1 "github.com/cmgsj/blob/pkg/proto/blob/v1"
)

func NewCommandList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List blobs",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			path := "/"

			if len(args) > 0 {
				path = args[0]
			}

			conn, err := newClientConn()
			if err != nil {
				return err
			}

			blobServiceClient := blobv1.NewBlobServiceClient(conn)

			request := &blobv1.ListBlobsRequest{}

			request.SetPath(path)

			response, err := blobServiceClient.ListBlobs(ctx, request)
			if err != nil {
				return err
			}

			for _, name := range response.GetNames() {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), name)
			}

			return nil
		},
	}

	return cmd
}

func NewCommandGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get blob",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			name := args[0]

			conn, err := newClientConn()
			if err != nil {
				return err
			}

			blobServiceClient := blobv1.NewBlobServiceClient(conn)

			request := &blobv1.GetBlobRequest{}

			request.SetName(name)

			response, err := blobServiceClient.GetBlob(ctx, request)
			if err != nil {
				return err
			}

			_, err = io.Copy(os.Stdout, bytes.NewReader(response.GetBlob().GetContent()))
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func NewCommandPut() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set blob",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			name := args[0]

			file := viper.GetString("file")

			var (
				content []byte
				err     error
			)

			if file == "" || file == "-" {
				content, err = io.ReadAll(os.Stdin)
			} else {
				content, err = os.ReadFile(file)
			}

			if err != nil {
				return err
			}

			conn, err := newClientConn()
			if err != nil {
				return err
			}

			blobServiceClient := blobv1.NewBlobServiceClient(conn)

			request := &blobv1.SetBlobRequest{}

			request.SetName(name)
			request.SetContent(content)

			_, err = blobServiceClient.SetBlob(ctx, request)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringP("file", "f", "", "input file")

	return cmd
}

func NewCommandDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete blob",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			name := args[0]

			conn, err := newClientConn()
			if err != nil {
				return err
			}

			blobServiceClient := blobv1.NewBlobServiceClient(conn)

			request := &blobv1.DeleteBlobRequest{}

			request.SetName(name)

			_, err = blobServiceClient.DeleteBlob(ctx, request)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func newClientConn() (*grpc.ClientConn, error) {
	address := viper.GetString("address")

	var opts []grpc.DialOption

	perRPCCredentials := newClientPerRPCCredentials()

	if perRPCCredentials != nil {
		opts = append(opts, grpc.WithPerRPCCredentials(perRPCCredentials))
	}

	transportCredentials, err := newClientTransportCredentials()
	if err != nil {
		return nil, err
	}

	opts = append(opts, grpc.WithTransportCredentials(transportCredentials))

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
