package write

import (
	"io"
	"os"

	cmdutil "github.com/cmgsj/blob/pkg/cmd/util"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/spf13/cobra"
)

type WriteOptions struct {
	IOStreams cmdutil.IOStreams
	Request   *blobv1.WriteFileRequest
	File      string
}

func NewWriteOptions(streams cmdutil.IOStreams) *WriteOptions {
	return &WriteOptions{
		IOStreams: streams,
		Request: &blobv1.WriteFileRequest{
			IsDir: new(bool),
		},
	}
}

func NewCmdWrite(f cmdutil.Factory, streams cmdutil.IOStreams) *cobra.Command {
	o := NewWriteOptions(streams)
	cmd := &cobra.Command{
		Use:  "write",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			stderr := o.IOStreams.Err
			cmdutil.CheckErr(o.Complete(f, cmd, args), stderr)
			cmdutil.CheckErr(o.Validate(), stderr)
			cmdutil.CheckErr(o.Run(f, cmd), stderr)
		},
	}
	cmd.Flags().BoolVar(o.Request.IsDir, "dir", *o.Request.IsDir, "create dir")
	cmd.Flags().StringVarP(&o.File, "file", "f", o.File, "file to write")
	return cmd
}

func (o *WriteOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	o.Request.FileName = args[0]
	return nil
}

func (o *WriteOptions) Validate() error {
	return nil
}

func (o *WriteOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	var content []byte
	var err error
	if o.File == "" {
		content, err = io.ReadAll(o.IOStreams.In)
	} else {
		content, err = os.ReadFile(o.File)
	}
	if err != nil {
		return err
	}
	o.Request.Content = content
	resp, err := f.BlobServiceClient().WriteFile(cmd.Context(), o.Request)
	if err != nil {
		return err
	}
	err = cmdutil.PrintJSON(o.IOStreams.Out, resp)
	return err
}
