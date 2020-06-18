package main

import (
	"git-get/pkg"
	"git-get/pkg/cfg"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmd = &cobra.Command{
	Use:          "git-list",
	Short:        "git list",
	RunE:         run,
	Args:         cobra.NoArgs,
	Version:      cfg.Version(),
	SilenceUsage: true,
}

func init() {
	cmd.PersistentFlags().StringP(cfg.KeyReposRoot, "r", "", "repos root")
	cmd.PersistentFlags().StringP(cfg.KeyPrivateKey, "p", "", "SSH private key path")
	cmd.PersistentFlags().BoolP(cfg.KeyFetch, "f", false, "Fetch from remotes when listing repositories")
	cmd.PersistentFlags().StringP(cfg.KeyOutput, "o", cfg.DefOutput, "output format.")

	viper.BindPFlag(cfg.KeyReposRoot, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))
	viper.BindPFlag(cfg.KeyPrivateKey, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))
	viper.BindPFlag(cfg.KeyFetch, cmd.PersistentFlags().Lookup(cfg.KeyFetch))
	viper.BindPFlag(cfg.KeyOutput, cmd.PersistentFlags().Lookup(cfg.KeyOutput))
}

func run(cmd *cobra.Command, args []string) error {
	cfg.Init()

	config := &pkg.ListCfg{
		Fetch:      viper.GetBool(cfg.KeyFetch),
		Output:     viper.GetString(cfg.KeyOutput),
		PrivateKey: viper.GetString(cfg.KeyPrivateKey),
		Root:       viper.GetString(cfg.KeyReposRoot),
	}

	return pkg.List(config)
}

func main() {
	cmd.Execute()
}
