package server

import (
	"github.com/cmgsj/blob/pkg/cli"
	"github.com/cmgsj/blob/pkg/cmd/server/health"
	"github.com/cmgsj/blob/pkg/cmd/server/start"
	"github.com/spf13/cobra"
)

func NewCmdServer(f cli.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "blob server",
		Run:   cli.Help,
	}
	cmd.AddCommand(start.NewCmdStart(f))
	cmd.AddCommand(health.NewCmdHealth(f))
	return cmd
}
