package blob

import (
	"github.com/spf13/cobra"

	"github.com/cmgsj/blob/pkg/cli"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

func NewCmdList(c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list blobs",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			path := "/"

			if len(args) > 0 {
				path = args[0]
			}

			blobClient, err := c.BlobServiceClient()
			if err != nil {
				return err
			}

			resp, err := blobClient.ListBlobs(ctx, &blobv1.ListBlobsRequest{
				Path: path,
			})
			if err != nil {
				return err
			}

			if resp.GetCount() == 0 {
				return nil
			}

			return cli.JSON(resp)
		},
	}

	return cmd
}
