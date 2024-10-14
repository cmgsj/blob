package blob

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cmgsj/blob/pkg/cli"
)

var version = "dev"

func NewCmdBlob() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blob",
		Short: "blob CLI",
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version,
	}

	cmd.PersistentFlags().String("grpc-address", "127.0.0.1:9090", "blob server grpc address")
	cmd.PersistentFlags().String("http-address", "127.0.0.1:8080", "blob server http address")

	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	viper.SetEnvPrefix("blob")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.BindPFlags(cmd.PersistentFlags())

	c := cli.NewConfig(cli.ConfigOptions{
		GRPCAddress: viper.GetString("grpc-address"),
		HTTPAddress: viper.GetString("http-address"),
	})

	cmd.AddCommand(
		NewCmdDelete(c),
		NewCmdGet(c),
		NewCmdList(c),
		NewCmdPut(c),
		NewCmdServer(c),
	)

	return cmd
}
