package main

import (
	"git-get/pkg"
	"git-get/pkg/cfg"
	"git-get/pkg/git"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const example = `  git get grdl/git-get
  git get https://github.com/grdl/git-get.git
  git get git@github.com:grdl/git-get.git
  git get -d path/to/dump/file`

var cmd = &cobra.Command{
	Use:          "git get <REPO>",
	Short:        "Clone git repository into an automatically created directory tree based on the repo's URL.",
	Example:      example,
	RunE:         run,
	Args:         cobra.MaximumNArgs(1), // TODO: add custom validator
	Version:      cfg.Version(),
	SilenceUsage: true, // We don't want to show usage on legit errors (eg, wrong path, repo already existing etc.)
}

func init() {
	cmd.PersistentFlags().StringP(cfg.KeyBranch, "b", "", "Branch (or tag) to checkout after cloning.")
	cmd.PersistentFlags().StringP(cfg.KeyDefaultHost, "t", cfg.Defaults[cfg.KeyDefaultHost], "Host to use when <REPO> doesn't have a specified host.")
	cmd.PersistentFlags().StringP(cfg.KeyDump, "d", "", "Path to a dump file listing repos to clone. Ignored when <REPO> argument is used.")
	cmd.PersistentFlags().StringP(cfg.KeyReposRoot, "r", cfg.Defaults[cfg.KeyReposRoot], "Path to repos root where repositories are cloned.")
	cmd.PersistentFlags().BoolP("help", "h", false, "Print this help and exit.")
	cmd.PersistentFlags().BoolP("version", "v", false, "Print version and exit.")

	viper.BindPFlag(cfg.KeyBranch, cmd.PersistentFlags().Lookup(cfg.KeyBranch))
	viper.BindPFlag(cfg.KeyDefaultHost, cmd.PersistentFlags().Lookup(cfg.KeyDefaultHost))
	viper.BindPFlag(cfg.KeyDump, cmd.PersistentFlags().Lookup(cfg.KeyDump))
	viper.BindPFlag(cfg.KeyReposRoot, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))
}

func run(cmd *cobra.Command, args []string) error {
	cfg.Init(&git.ConfigGlobal{})

	var url string
	if len(args) > 0 {
		url = args[0]
	}

	config := &pkg.GetCfg{
		Branch:  viper.GetString(cfg.KeyBranch),
		DefHost: viper.GetString(cfg.KeyDefaultHost),
		Dump:    viper.GetString(cfg.KeyDump),
		Root:    viper.GetString(cfg.KeyReposRoot),
		URL:     url,
	}
	return pkg.Get(config)
}

func main() {
	cmd.Execute()
}
