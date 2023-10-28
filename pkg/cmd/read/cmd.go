package read

import (
	"context"
	"fmt"

	"github.com/cmgsj/blob/pkg/cli"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/spf13/cobra"
)

func NewCmdRead(f cli.Factory) *cobra.Command {
	o := NewReadOptions(f)
	cmd := &cobra.Command{
		Use:   "read",
		Short: "read blob",
		Args:  cobra.ExactArgs(1),
		Run:   cli.Run(o),
	}
	return cmd
}

type ReadOptions struct {
	cli.Factory
	Request *blobv1.GetBlobRequest
}

func NewReadOptions(f cli.Factory) *ReadOptions {
	return &ReadOptions{
		Factory: f,
		Request: &blobv1.GetBlobRequest{},
	}
}

func (o *ReadOptions) Complete(ctx context.Context, cmd *cobra.Command, args []string) error {
	o.Request.BlobName = args[0]
	return nil
}

func (o *ReadOptions) Validate(ctx context.Context) error {
	return nil
}

func (o *ReadOptions) Run(ctx context.Context) error {
	client, err := o.BlobServiceClient(ctx)
	if err != nil {
		return err
	}
	resp, err := client.GetBlob(ctx, o.Request)
	if err != nil {
		return err
	}
	_, err = fmt.Printf("%s\n", resp.GetBlob().GetContent())
	return err
}
