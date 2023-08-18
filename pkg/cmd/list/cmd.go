package list

import (
	cmdutil "github.com/cmgsj/blob/pkg/cmd/util"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/spf13/cobra"
)

type ListOptions struct {
	IOStreams cmdutil.IOStreams
	Request   *blobv1.ListFilesRequest
}

func NewListOptions(streams cmdutil.IOStreams) *ListOptions {
	return &ListOptions{
		IOStreams: streams,
		Request: &blobv1.ListFilesRequest{
			Path: "/",
		},
	}
}

func NewCmdList(f cmdutil.Factory, streams cmdutil.IOStreams) *cobra.Command {
	o := NewListOptions(streams)
	cmd := &cobra.Command{
		Use:  "list",
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			stderr := o.IOStreams.Err
			cmdutil.CheckErr(o.Complete(f, cmd, args), stderr)
			cmdutil.CheckErr(o.Validate(), stderr)
			cmdutil.CheckErr(o.Run(f, cmd), stderr)
		},
	}
	return cmd
}

func (o *ListOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		o.Request.Path = args[0]
	}
	return nil
}

func (o *ListOptions) Validate() error {
	return nil
}

func (o *ListOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	resp, err := f.BlobServiceClient().ListFiles(cmd.Context(), o.Request)
	if err != nil {
		return err
	}
	err = cmdutil.PrintJSON(o.IOStreams.Out, resp)
	return err
}
