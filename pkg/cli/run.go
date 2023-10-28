package cli

import (
	"context"

	"github.com/spf13/cobra"
)

type Command interface {
	Factory
	Complete(ctx context.Context, cmd *cobra.Command, args []string) error
	Validate(ctx context.Context) error
	Run(ctx context.Context) error
}

func Run(c Command) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		Check(c.Complete(ctx, cmd, args))
		Check(c.Validate(ctx))
		Check(c.Run(ctx))
	}
}

func Help(cmd *cobra.Command, args []string) {
	cmd.Help()
}
