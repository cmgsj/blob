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
			req := &blobv1.RemoveBlobRequest{
				BlobName: args[0],
			}

			client, err := c.BlobServiceClient()
			if err != nil {
				return err
			}

			_, err = client.RemoveBlob(cmd.Context(), req)
			return err
		},
	}

	return cmd
}
