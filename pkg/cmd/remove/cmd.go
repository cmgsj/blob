package remove

import (
	cmdutil "github.com/cmgsj/blob/pkg/cmd/util"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/spf13/cobra"
)

type RemoveOptions struct {
	IOStreams cmdutil.IOStreams
	Request   *blobv1.RemoveBlobRequest
}

func NewRemoveOptions(streams cmdutil.IOStreams) *RemoveOptions {
	return &RemoveOptions{
		IOStreams: streams,
		Request:   &blobv1.RemoveBlobRequest{},
	}
}

func NewCmdRemove(f cmdutil.Factory, streams cmdutil.IOStreams) *cobra.Command {
	o := NewRemoveOptions(streams)
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Short:   "remove blob",
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

func (o *RemoveOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	o.Request.BlobName = args[0]
	return nil
}

func (o *RemoveOptions) Validate() error {
	return nil
}

func (o *RemoveOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	_, err := f.BlobServiceClient().RemoveBlob(cmd.Context(), o.Request)
	return err
}
