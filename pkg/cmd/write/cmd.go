package write

import (
	"context"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/cmgsj/blob/pkg/cli"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

func NewCmdWrite(f cli.Factory) *cobra.Command {
	o := NewWriteOptions(f)
	cmd := &cobra.Command{
		Use:   "write",
		Short: "write blob",
		Args:  cobra.ExactArgs(1),
		Run:   cli.Run(o),
	}
	cmd.Flags().StringVarP(&o.File, "file", "f", o.File, "input file")
	return cmd
}

type WriteOptions struct {
	cli.Factory
	Request *blobv1.WriteBlobRequest
	File    string
}

func NewWriteOptions(f cli.Factory) *WriteOptions {
	return &WriteOptions{
		Factory: f,
		Request: &blobv1.WriteBlobRequest{},
	}
}

func (o *WriteOptions) Complete(ctx context.Context, cmd *cobra.Command, args []string) error {
	o.Request.BlobName = args[0]
	return nil
}

func (o *WriteOptions) Validate(ctx context.Context) error {
	return nil
}

func (o *WriteOptions) Run(ctx context.Context) error {
	var content []byte
	var err error
	if o.File == "" {
		content, err = io.ReadAll(os.Stdin)
	} else {
		content, err = os.ReadFile(o.File)
	}
	if err != nil {
		return err
	}
	o.Request.Content = content
	client, err := o.BlobServiceClient(ctx)
	if err != nil {
		return err
	}
	_, err = client.WriteBlob(ctx, o.Request)
	return err
}
