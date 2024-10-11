package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cmgsj/blob/pkg/cli"
	"github.com/cmgsj/blob/pkg/cmd/get"
	"github.com/cmgsj/blob/pkg/cmd/list"
	"github.com/cmgsj/blob/pkg/cmd/read"
	"github.com/cmgsj/blob/pkg/cmd/remove"
	"github.com/cmgsj/blob/pkg/cmd/server"
	"github.com/cmgsj/blob/pkg/cmd/write"
)

var version = "dev"

func NewCmdBlob() *cobra.Command {
	o := NewBlobOptions()
	cmd := &cobra.Command{
		Use:     "blob",
		Short:   "blob CLI",
		Run:     cli.Help,
		Version: version,
	}
	f := cli.NewFactory(o)
	cmd.AddCommand(get.NewCmdGet(f))
	cmd.AddCommand(list.NewCmdList(f))
	cmd.AddCommand(read.NewCmdRead(f))
	cmd.AddCommand(remove.NewCmdRemove(f))
	cmd.AddCommand(write.NewCmdWrite(f))
	cmd.AddCommand(server.NewCmdServer(f))
	cmd.PersistentFlags().StringVar(&o.GRPCAddress, "grpc-address", o.GRPCAddress, "blob server grpc address")
	cmd.PersistentFlags().StringVar(&o.HTTPAddress, "http-address", o.HTTPAddress, "blob server http address")
	return cmd
}

type BlobOptions struct {
	GRPCAddress string
	HTTPAddress string
}

func NewBlobOptions() *BlobOptions {
	return &BlobOptions{
		GRPCAddress: "127.0.0.1:9090",
		HTTPAddress: "127.0.0.1:8080",
	}
}

func (o *BlobOptions) GetGRPCAddress() string {
	return o.GRPCAddress
}

func (o *BlobOptions) GetHTTPAddress() string {
	return o.HTTPAddress
}
