package blob

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	apiv1 "github.com/cmgsj/blob/pkg/proto/blob/api/v1"
	"github.com/cmgsj/blob/pkg/proto/blob/api/v1/apiv1connect"
)

func NewCommandGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get blob",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			name := args[0]

			baseURL := viper.GetString("base-url")

			blobClient := apiv1connect.NewBlobServiceClient(http.DefaultClient, baseURL)

			request := &apiv1.GetBlobRequest{}

			request.SetName(name)

			response, err := blobClient.GetBlob(ctx, connect.NewRequest(request))
			if err != nil {
				return err
			}

			_, err = io.Copy(os.Stdout, bytes.NewReader(response.Msg.GetBlob().GetContent()))
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String("base-url", "http://localhost:8080", "blob service base url")

	return cmd
}
