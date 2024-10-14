package blob

import (
	"github.com/spf13/cobra"

	"github.com/cmgsj/blob/pkg/cli"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

func NewCmdDelete(c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete blob",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			blobName := args[0]

			blobClient, err := c.BlobServiceClient()
			if err != nil {
				return err
			}

			_, err = blobClient.DeleteBlob(ctx, &blobv1.DeleteBlobRequest{
				BlobName: blobName,
			})
			return err
		},
	}

	return cmd
}
