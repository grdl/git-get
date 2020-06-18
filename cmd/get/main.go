package main

import (
	"git-get/pkg"
	"git-get/pkg/cfg"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmd = &cobra.Command{
	Use:          "git-get <repo>",
	Short:        "git get",
	RunE:         run,
	Args:         cobra.MaximumNArgs(1), // TODO: add custom validator
	Version:      cfg.Version(),
	SilenceUsage: true,
}

func init() {
	cmd.PersistentFlags().StringP(cfg.KeyReposRoot, "r", "", "repos root")
	cmd.PersistentFlags().StringP(cfg.KeyPrivateKey, "p", "", "SSH private key path")
	cmd.PersistentFlags().StringP(cfg.KeyDump, "d", "", "Dump file path")
	cmd.PersistentFlags().StringP(cfg.KeyBranch, "b", cfg.DefBranch, "Branch (or tag) to checkout after cloning")

	viper.BindPFlag(cfg.KeyReposRoot, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))
	viper.BindPFlag(cfg.KeyPrivateKey, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))
	viper.BindPFlag(cfg.KeyDump, cmd.PersistentFlags().Lookup(cfg.KeyDump))
	viper.BindPFlag(cfg.KeyBranch, cmd.PersistentFlags().Lookup(cfg.KeyBranch))
}

func run(cmd *cobra.Command, args []string) error {
	cfg.Init()

	var url string
	if len(args) > 0 {
		url = args[0]
	}

	config := &pkg.GetCfg{
		Branch: viper.GetString(cfg.KeyBranch),
		Dump:   viper.GetString(cfg.KeyDump),
		Root:   viper.GetString(cfg.KeyReposRoot),
		URL:    url,
	}
	return pkg.Get(config)
}

func main() {
	cmd.Execute()
}
