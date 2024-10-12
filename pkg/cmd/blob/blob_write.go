package blob

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/cmgsj/blob/pkg/cli"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

func NewCmdWrite(c *cli.Config) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "write",
		Short: "write blob",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			req := &blobv1.WriteBlobRequest{
				BlobName: args[0],
			}

			var err error

			if file == "" {
				req.Content, err = io.ReadAll(os.Stdin)
			} else {
				req.Content, err = os.ReadFile(file)
			}
			if err != nil {
				return err
			}

			blobClient, err := c.BlobServiceClient()
			if err != nil {
				return err
			}

			_, err = blobClient.WriteBlob(ctx, req)
			return err
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", file, "input file")

	return cmd
}
