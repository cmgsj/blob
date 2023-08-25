package read

import (
	"fmt"

	cmdutil "github.com/cmgsj/blob/pkg/cmd/util"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/spf13/cobra"
)

type ReadOptions struct {
	IOStreams cmdutil.IOStreams
	Request   *blobv1.GetBlobRequest
}

func NewReadOptions(streams cmdutil.IOStreams) *ReadOptions {
	return &ReadOptions{
		IOStreams: streams,
		Request:   &blobv1.GetBlobRequest{},
	}
}

func NewCmdRead(f cmdutil.Factory, streams cmdutil.IOStreams) *cobra.Command {
	o := NewReadOptions(streams)
	cmd := &cobra.Command{
		Use:     "read",
		Aliases: []string{"r"},
		Short:   "read blob",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			stderr := o.IOStreams.Err
			cmdutil.CheckErr(o.Complete(f, cmd, args), stderr)
			cmdutil.CheckErr(o.Validate(), stderr)
			cmdutil.CheckErr(o.Run(f, cmd), stderr)
		},
	}
	return cmd
}

func (o *ReadOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	o.Request.BlobName = args[0]
	return nil
}

func (o *ReadOptions) Validate() error {
	return nil
}

func (o *ReadOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	resp, err := f.BlobServiceClient().GetBlob(cmd.Context(), o.Request)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(o.IOStreams.Out, "%s\n", resp.GetBlob().GetContent())
	return err
}
