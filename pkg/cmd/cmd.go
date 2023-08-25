package cmd

import (
	"context"

	"github.com/cmgsj/blob/pkg/cmd/get"
	"github.com/cmgsj/blob/pkg/cmd/list"
	"github.com/cmgsj/blob/pkg/cmd/read"
	"github.com/cmgsj/blob/pkg/cmd/remove"
	"github.com/cmgsj/blob/pkg/cmd/server"
	cmdutil "github.com/cmgsj/blob/pkg/cmd/util"
	"github.com/cmgsj/blob/pkg/cmd/write"
	"github.com/cmgsj/blob/pkg/version"
	"github.com/spf13/cobra"
)

type BlobOptions struct {
	IOStreams   cmdutil.IOStreams
	HTTPAddress string
	GRPCAddress string
}

func NewBlobOptions(streams cmdutil.IOStreams) *BlobOptions {
	return &BlobOptions{
		IOStreams:   streams,
		HTTPAddress: "127.0.0.1:8080",
		GRPCAddress: "127.0.0.1:9090",
	}
}

func NewCmdBlob() *cobra.Command {
	o := NewBlobOptions(cmdutil.NewOSStreams())
	cmd := &cobra.Command{
		Use:     "blob",
		Short:   "blob CLI",
		Run:     cmdutil.RunHelp,
		Version: version.Version,
	}
	ctx := context.Background()
	cmd.SetContext(ctx)
	streams := o.IOStreams
	stderr := streams.Err
	f, err := cmdutil.NewFactory(cmd.Context(), o.GRPCAddress)
	cmdutil.CheckErr(err, stderr)
	cmd.AddCommand(get.NewCmdGet(f, streams))
	cmd.AddCommand(list.NewCmdList(f, streams))
	cmd.AddCommand(read.NewCmdRead(f, streams))
	cmd.AddCommand(remove.NewCmdRemove(f, streams))
	cmd.AddCommand(write.NewCmdWrite(f, streams))
	cmd.AddCommand(server.NewCmdServer(f, streams))
	cmd.PersistentFlags().StringVar(&o.HTTPAddress, "http-address", o.HTTPAddress, "blob server http address")
	cmd.PersistentFlags().StringVar(&o.GRPCAddress, "grpc-address", o.GRPCAddress, "blob server grpc address")
	return cmd
}
