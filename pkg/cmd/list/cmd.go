package list

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cmgsj/blob/pkg/cli"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

func NewCmdList(f cli.Factory) *cobra.Command {
	o := NewListOptions(f)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list blobs",
		Args:  cobra.MaximumNArgs(1),
		Run:   cli.Run(o),
	}
	return cmd
}

type ListOptions struct {
	cli.Factory
	Request *blobv1.ListBlobsRequest
}

func NewListOptions(f cli.Factory) *ListOptions {
	return &ListOptions{
		Factory: f,
		Request: &blobv1.ListBlobsRequest{
			Path: "/",
		},
	}
}

func (o *ListOptions) Complete(ctx context.Context, cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		o.Request.Path = args[0]
	}
	return nil
}

func (o *ListOptions) Validate(ctx context.Context) error {
	return nil
}

func (o *ListOptions) Run(ctx context.Context) error {
	client, err := o.BlobServiceClient(ctx)
	if err != nil {
		return err
	}
	resp, err := client.ListBlobs(ctx, o.Request)
	if err != nil {
		return err
	}
	if resp.GetCount() == 0 {
		return nil
	}
	return cli.JSON(resp)
}
