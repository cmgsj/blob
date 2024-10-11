package remove

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cmgsj/blob/pkg/cli"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

func NewCmdRemove(f cli.Factory) *cobra.Command {
	o := NewRemoveOptions(f)
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "remove blob",
		Args:  cobra.ExactArgs(1),
		Run:   cli.Run(o),
	}
	return cmd
}

type RemoveOptions struct {
	cli.Factory
	Request *blobv1.RemoveBlobRequest
}

func NewRemoveOptions(f cli.Factory) *RemoveOptions {
	return &RemoveOptions{
		Factory: f,
		Request: &blobv1.RemoveBlobRequest{},
	}
}

func (o *RemoveOptions) Complete(ctx context.Context, cmd *cobra.Command, args []string) error {
	o.Request.BlobName = args[0]
	return nil
}

func (o *RemoveOptions) Validate(ctx context.Context) error {
	return nil
}

func (o *RemoveOptions) Run(ctx context.Context) error {
	client, err := o.BlobServiceClient(ctx)
	if err != nil {
		return err
	}
	_, err = client.RemoveBlob(ctx, o.Request)
	return err
}
