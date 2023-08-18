package util

import "github.com/spf13/cobra"

func RunHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
