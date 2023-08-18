package server

import (
	"github.com/cmgsj/blob/pkg/cmd/server/start"
	cmdutil "github.com/cmgsj/blob/pkg/cmd/util"
	"github.com/spf13/cobra"
)

type ServerOptions struct {
	IOStreams cmdutil.IOStreams
}

func NewServerOptions(streams cmdutil.IOStreams) *ServerOptions {
	return &ServerOptions{
		IOStreams: streams,
	}
}

func NewCmdServer(streams cmdutil.IOStreams) *cobra.Command {
	_ = NewServerOptions(streams)
	cmd := &cobra.Command{
		Use: "server",
		Run: cmdutil.RunHelp,
	}
	cmd.AddCommand(start.NewCmdStart(streams))
	return cmd
}
