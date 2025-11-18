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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			viper.AutomaticEnv()
			viper.AllowEmptyEnv(true)
			viper.SetEnvPrefix(cmd.Name())
			viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			return viper.BindPFlags(cmd.PersistentFlags())
		},
	}

	cmd.PersistentFlags().String("address", "localhost:2562", "server address")

	cmd.AddCommand(
		NewCommandList(),
		NewCommandGet(),
		NewCommandPut(),
		NewCommandDelete(),
		NewCommandHealth(),
		NewCommandServer(),
	)

	return cmd
}
