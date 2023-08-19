package delete

import (
	cmdutil "github.com/cmgsj/blob/pkg/cmd/util"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/spf13/cobra"
)

type DeleteOptions struct {
	IOStreams cmdutil.IOStreams
	Request   *blobv1.DeleteFileRequest
}

func NewDeleteOptions(streams cmdutil.IOStreams) *DeleteOptions {
	return &DeleteOptions{
		IOStreams: streams,
		Request:   &blobv1.DeleteFileRequest{},
	}
}

func NewCmdDelete(f cmdutil.Factory, streams cmdutil.IOStreams) *cobra.Command {
	o := NewDeleteOptions(streams)
	cmd := &cobra.Command{
		Use:  "delete",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			stderr := o.IOStreams.Err
			cmdutil.CheckErr(o.Complete(f, cmd, args), stderr)
			cmdutil.CheckErr(o.Validate(), stderr)
			cmdutil.CheckErr(o.Run(f, cmd), stderr)
		},
	}
	return cmd
}

func (o *DeleteOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	o.Request.FileName = args[0]
	return nil
}

func (o *DeleteOptions) Validate() error {
	return nil
}

func (o *DeleteOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	_, err := f.BlobServiceClient().DeleteFile(cmd.Context(), o.Request)
	return err
}
