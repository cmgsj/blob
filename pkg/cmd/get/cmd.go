package get

import (
	cmdutil "github.com/cmgsj/blob/pkg/cmd/util"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/spf13/cobra"
)

type GetOptions struct {
	IOStreams cmdutil.IOStreams
	Request   *blobv1.GetFileRequest
}

func NewGetOptions(streams cmdutil.IOStreams) *GetOptions {
	return &GetOptions{
		IOStreams: streams,
		Request:   &blobv1.GetFileRequest{},
	}
}

func NewCmdGet(f cmdutil.Factory, streams cmdutil.IOStreams) *cobra.Command {
	o := NewGetOptions(streams)
	cmd := &cobra.Command{
		Use:  "get",
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

func (o *GetOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	o.Request.FileName = args[0]
	return nil
}

func (o *GetOptions) Validate() error {
	return nil
}

func (o *GetOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	resp, err := f.BlobServiceClient().GetFile(cmd.Context(), o.Request)
	if err != nil {
		return err
	}
	err = cmdutil.PrintJSON(o.IOStreams.Out, resp)
	return err
}
