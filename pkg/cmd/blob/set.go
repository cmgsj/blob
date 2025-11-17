package blob

import (
	"io"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	apiv1 "github.com/cmgsj/blob/pkg/proto/blob/api/v1"
	"github.com/cmgsj/blob/pkg/proto/blob/api/v1/apiv1connect"
)

func NewCommandPut() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set blob",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			name := args[0]

			baseURL := viper.GetString("base-url")
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

			blobClient := apiv1connect.NewBlobServiceClient(http.DefaultClient, baseURL)

			request := &apiv1.SetBlobRequest{}

			request.SetName(name)
			request.SetContent(content)

			_, err = blobClient.SetBlob(ctx, connect.NewRequest(request))
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String("base-url", "http://localhost:8080", "blob service base url")
	cmd.Flags().StringP("file", "f", "", "input file")

	_ = viper.BindPFlags(cmd.Flags())

	return cmd
}
