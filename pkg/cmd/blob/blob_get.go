package blob

import (
	"github.com/spf13/cobra"

	"github.com/cmgsj/blob/pkg/cli"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

func NewCmdGet(c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get blob",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			req := &blobv1.GetBlobRequest{
				BlobName: args[0],
			}

			client, err := c.BlobServiceClient()
			if err != nil {
				return err
			}

			resp, err := client.GetBlob(cmd.Context(), req)
			if err != nil {
				return err
			}

			resp.Blob.Content = nil

			return cli.JSON(resp)
		},
	}

	return cmd
}
