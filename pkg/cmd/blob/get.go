package blob

import (
	"bytes"
	"io"
	"os"

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
			ctx := cmd.Context()

			blobName := args[0]

			blobClient, err := c.BlobServiceClient()
			if err != nil {
				return err
			}

			resp, err := blobClient.GetBlob(ctx, &blobv1.GetBlobRequest{
				BlobName: blobName,
			})
			if err != nil {
				return err
			}

			_, err = io.Copy(os.Stdout, bytes.NewReader(resp.GetBlob().GetContent()))
			return err
		},
	}

	return cmd
}
