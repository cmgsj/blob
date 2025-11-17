package blob

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "dev"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blob",
		Short: "Blob CLI",
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version,
	}

	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	viper.SetEnvPrefix(cmd.Name())
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	cmd.AddCommand(
		NewCommandServer(),
		NewCommandHealth(),
		NewCommandList(),
		NewCommandGet(),
		NewCommandPut(),
		NewCommandDelete(),
	)

	return cmd
}
