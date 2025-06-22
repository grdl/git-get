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
  git get grdl/git-get --depth 1
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
	cmd.PersistentFlags().IntP(cfg.KeyDepth, "D", 0, "Create a shallow clone with a history truncated to the specified number of commits.")
	cmd.PersistentFlags().StringP(cfg.KeyDefaultHost, "t", cfg.Defaults[cfg.KeyDefaultHost], "Host to use when <REPO> doesn't have a specified host.")
	cmd.PersistentFlags().StringP(cfg.KeyDefaultScheme, "c", cfg.Defaults[cfg.KeyDefaultScheme], "Scheme to use when <REPO> doesn't have a specified scheme.")
	cmd.PersistentFlags().StringP(cfg.KeyDump, "d", "", "Path to a dump file listing repos to clone. Ignored when <REPO> argument is used.")
	cmd.PersistentFlags().BoolP(cfg.KeySkipHost, "s", false, "Don't create a directory for host.")
	cmd.PersistentFlags().StringP(cfg.KeyReposRoot, "r", cfg.Defaults[cfg.KeyReposRoot], "Path to repos root where repositories are cloned.")
	cmd.PersistentFlags().BoolP("help", "h", false, "Print this help and exit.")
	cmd.PersistentFlags().BoolP("version", "v", false, "Print version and exit.")

	viper.BindPFlag(cfg.KeyBranch, cmd.PersistentFlags().Lookup(cfg.KeyBranch))
	viper.BindPFlag(cfg.KeyDepth, cmd.PersistentFlags().Lookup(cfg.KeyDepth))
	viper.BindPFlag(cfg.KeyDefaultHost, cmd.PersistentFlags().Lookup(cfg.KeyDefaultHost))
	viper.BindPFlag(cfg.KeyDefaultScheme, cmd.PersistentFlags().Lookup(cfg.KeyDefaultScheme))
	viper.BindPFlag(cfg.KeyDump, cmd.PersistentFlags().Lookup(cfg.KeyDump))
	viper.BindPFlag(cfg.KeyReposRoot, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))
	viper.BindPFlag(cfg.KeySkipHost, cmd.PersistentFlags().Lookup(cfg.KeySkipHost))

	cfg.Init(&git.ConfigGlobal{})
}

func run(cmd *cobra.Command, args []string) error {
	var url string
	if len(args) > 0 {
		url = args[0]
	}

	cfg.Expand(cfg.KeyReposRoot)

	config := &pkg.GetCfg{
		Branch:    viper.GetString(cfg.KeyBranch),
		DefHost:   viper.GetString(cfg.KeyDefaultHost),
		DefScheme: viper.GetString(cfg.KeyDefaultScheme),
		Depth:     viper.GetInt(cfg.KeyDepth),
		Dump:      viper.GetString(cfg.KeyDump),
		SkipHost:  viper.GetBool(cfg.KeySkipHost),
		Root:      viper.GetString(cfg.KeyReposRoot),
		URL:       url,
	}
	return pkg.Get(config)
}

func main() {
	cmd.Execute()
}
