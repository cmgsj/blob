package get

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cmgsj/blob/pkg/cli"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

func NewCmdGet(f cli.Factory) *cobra.Command {
	o := NewGetOptions(f)
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get blob",
		Args:  cobra.ExactArgs(1),
		Run:   cli.Run(o),
	}
	return cmd
}

type GetOptions struct {
	cli.Factory
	Request *blobv1.GetBlobRequest
}

func NewGetOptions(f cli.Factory) *GetOptions {
	return &GetOptions{
		Factory: f,
		Request: &blobv1.GetBlobRequest{},
	}
}

func (o *GetOptions) Complete(ctx context.Context, cmd *cobra.Command, args []string) error {
	o.Request.BlobName = args[0]
	return nil
}

func (o *GetOptions) Validate(ctx context.Context) error {
	return nil
}

func (o *GetOptions) Run(ctx context.Context) error {
	client, err := o.BlobServiceClient(ctx)
	if err != nil {
		return err
	}
	resp, err := client.GetBlob(ctx, o.Request)
	if err != nil {
		return err
	}
	resp.Blob.Content = nil
	return cli.JSON(resp)
}
