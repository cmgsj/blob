package rename

import (
	cmdutil "github.com/cmgsj/blob/pkg/cmd/util"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/spf13/cobra"
)

type RenameOptions struct {
	IOStreams cmdutil.IOStreams
	Request   *blobv1.RenameFileRequest
}

func NewRenameOptions(streams cmdutil.IOStreams) *RenameOptions {
	return &RenameOptions{
		IOStreams: streams,
		Request:   &blobv1.RenameFileRequest{},
	}
}

func NewCmdRename(f cmdutil.Factory, streams cmdutil.IOStreams) *cobra.Command {
	o := NewRenameOptions(streams)
	cmd := &cobra.Command{
		Use:  "rename",
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			stderr := o.IOStreams.Err
			cmdutil.CheckErr(o.Complete(f, cmd, args), stderr)
			cmdutil.CheckErr(o.Validate(), stderr)
			cmdutil.CheckErr(o.Run(f, cmd), stderr)
		},
	}
	return cmd
}

func (o *RenameOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	o.Request.FileName = args[0]
	o.Request.NewFileName = args[1]
	return nil
}

func (o *RenameOptions) Validate() error {
	return nil
}

func (o *RenameOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	_, err := f.BlobServiceClient().RenameFile(cmd.Context(), o.Request)
	return err
}
