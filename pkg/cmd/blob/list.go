package blob

import (
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	apiv1 "github.com/cmgsj/blob/pkg/proto/blob/api/v1"
	"github.com/cmgsj/blob/pkg/proto/blob/api/v1/apiv1connect"
)

func NewCommandList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List blobs",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			path := "/"

			if len(args) > 0 {
				path = args[0]
			}

			baseURL := viper.GetString("base-url")

			blobClient := apiv1connect.NewBlobServiceClient(http.DefaultClient, baseURL)

			request := &apiv1.ListBlobsRequest{}

			request.SetPath(path)

			response, err := blobClient.ListBlobs(ctx, connect.NewRequest(request))
			if err != nil {
				return err
			}

			for _, name := range response.Msg.GetNames() {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), name)
			}

			return nil
		},
	}

	cmd.Flags().String("base-url", "http://localhost:8080", "blob service base url")

	return cmd
}
