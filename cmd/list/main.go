package main

import (
	"fmt"
	"git-get/pkg"
	"git-get/pkg/cfg"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmd = &cobra.Command{
	Use:          "git list",
	Short:        "List all repositories cloned by 'git get' and their status.",
	RunE:         run,
	Args:         cobra.NoArgs,
	Version:      cfg.Version(),
	SilenceUsage: true, // We don't want to show usage on legit errors (eg, wrong path, repo already existing etc.)
}

func init() {
	cmd.PersistentFlags().BoolP(cfg.KeyFetch, "f", false, "First fetch from remotes before listing repositories.")
	cmd.PersistentFlags().StringP(cfg.KeyOutput, "o", cfg.DefOutput, fmt.Sprintf("Output format. Allowed values: [%s].", strings.Join(cfg.AllowedOut, ", ")))
	cmd.PersistentFlags().StringP(cfg.KeyPrivateKey, "p", "", "Path to SSH private key. (default \"~/.ssh/id_rsa\")")
	cmd.PersistentFlags().StringP(cfg.KeyReposRoot, "r", "", "Path to repos root where repositories are cloned. (default \"~/repositories\")")
	cmd.PersistentFlags().BoolP("help", "h", false, "Print this help and exit.")
	cmd.PersistentFlags().BoolP("version", "v", false, "Print version and exit.")

	viper.BindPFlag(cfg.KeyFetch, cmd.PersistentFlags().Lookup(cfg.KeyFetch))
	viper.BindPFlag(cfg.KeyOutput, cmd.PersistentFlags().Lookup(cfg.KeyOutput))
	viper.BindPFlag(cfg.KeyPrivateKey, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))
	viper.BindPFlag(cfg.KeyReposRoot, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))

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
