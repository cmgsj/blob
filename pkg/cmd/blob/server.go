package blob

import (
	"github.com/spf13/cobra"

	"github.com/cmgsj/blob/pkg/cli"
)

func NewCmdServer(c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "blob server",
	}

	cmd.AddCommand(
		NewCmdServerHealth(c),
		NewCmdServerStart(c),
	)

	return cmd
}
