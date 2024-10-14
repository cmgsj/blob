package blob

import (
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cmgsj/blob/pkg/cli"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

func NewCmdPut(c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "put",
		Short: "put blob",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			blobName := args[0]
			file := viper.GetString("file")

			var content []byte
			var err error

			if file == "" || file == "-" {
				content, err = io.ReadAll(os.Stdin)
			} else {
				content, err = os.ReadFile(file)
			}
			if err != nil {
				return err
			}

			blobClient, err := c.BlobServiceClient()
			if err != nil {
				return err
			}

			_, err = blobClient.PutBlob(ctx, &blobv1.PutBlobRequest{
				BlobName: blobName,
				Content:  content,
			})
			return err
		},
	}

	cmd.Flags().StringP("file", "f", "", "input file")

	viper.BindPFlags(cmd.Flags())

	return cmd
}
