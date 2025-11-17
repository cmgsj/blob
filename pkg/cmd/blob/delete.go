package blob

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	apiv1 "github.com/cmgsj/blob/pkg/proto/blob/api/v1"
	"github.com/cmgsj/blob/pkg/proto/blob/api/v1/apiv1connect"
)

func NewCommandDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete blob",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			name := args[0]

			baseURL := viper.GetString("base-url")

			blobClient := apiv1connect.NewBlobServiceClient(http.DefaultClient, baseURL)

			request := &apiv1.DeleteBlobRequest{}

			request.SetName(name)

			_, err := blobClient.DeleteBlob(ctx, connect.NewRequest(request))
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String("base-url", "http://localhost:8080", "blob service base url")

	return cmd
}
