package blob

import (
	"github.com/spf13/cobra"

	"github.com/cmgsj/blob/pkg/cli"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

func NewCmdRemove(c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "remove blob",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			blobName := args[0]

			blobClient, err := c.BlobServiceClient()
			if err != nil {
				return err
			}

			_, err = blobClient.RemoveBlob(ctx, &blobv1.RemoveBlobRequest{
				BlobName: blobName,
			})
			return err
		},
	}

	return cmd
}
