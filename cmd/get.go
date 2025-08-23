package main

import (
	"git-get/pkg"
	"git-get/pkg/cfg"
	"git-get/pkg/git"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const getExample = `  git get grdl/git-get
  git get https://github.com/grdl/git-get.git
  git get git@github.com:grdl/git-get.git
  git get -d path/to/dump/file`

func newGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "git get <REPO>",
		Short:        "Clone git repository into an automatically created directory tree based on the repo's URL.",
		Example:      getExample,
		RunE:         runGetCommand,
		Args:         cobra.MaximumNArgs(1), // TODO: add custom validator
		Version:      cfg.Version(),
		SilenceUsage: true, // We don't want to show usage on legit errors (eg, wrong path, repo already existing etc.)
	}

	cmd.PersistentFlags().StringP(cfg.KeyBranch, "b", "", "Branch (or tag) to checkout after cloning.")
	cmd.PersistentFlags().StringP(cfg.KeyDefaultHost, "t", cfg.Defaults[cfg.KeyDefaultHost], "Host to use when <REPO> doesn't have a specified host.")
	cmd.PersistentFlags().StringP(cfg.KeyDefaultScheme, "c", cfg.Defaults[cfg.KeyDefaultScheme], "Scheme to use when <REPO> doesn't have a specified scheme.")
	cmd.PersistentFlags().StringP(cfg.KeyDump, "d", "", "Path to a dump file listing repos to clone. Ignored when <REPO> argument is used.")
	cmd.PersistentFlags().BoolP(cfg.KeySkipHost, "s", false, "Don't create a directory for host.")
	cmd.PersistentFlags().StringP(cfg.KeyReposRoot, "r", cfg.Defaults[cfg.KeyReposRoot], "Path to repos root where repositories are cloned.")
	cmd.PersistentFlags().BoolP("help", "h", false, "Print this help and exit.")
	cmd.PersistentFlags().BoolP("version", "v", false, "Print version and exit.")

	viper.BindPFlag(cfg.KeyBranch, cmd.PersistentFlags().Lookup(cfg.KeyBranch))
	viper.BindPFlag(cfg.KeyDefaultHost, cmd.PersistentFlags().Lookup(cfg.KeyDefaultHost))
	viper.BindPFlag(cfg.KeyDefaultScheme, cmd.PersistentFlags().Lookup(cfg.KeyDefaultScheme))
	viper.BindPFlag(cfg.KeyDump, cmd.PersistentFlags().Lookup(cfg.KeyDump))
	viper.BindPFlag(cfg.KeyReposRoot, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))
	viper.BindPFlag(cfg.KeySkipHost, cmd.PersistentFlags().Lookup(cfg.KeySkipHost))

	return cmd
}

func runGetCommand(cmd *cobra.Command, args []string) error {
	var url string
	if len(args) > 0 {
		url = args[0]
	}

	cfg.Expand(cfg.KeyReposRoot)

	config := &pkg.GetCfg{
		Branch:    viper.GetString(cfg.KeyBranch),
		DefHost:   viper.GetString(cfg.KeyDefaultHost),
		DefScheme: viper.GetString(cfg.KeyDefaultScheme),
		Dump:      viper.GetString(cfg.KeyDump),
		SkipHost:  viper.GetBool(cfg.KeySkipHost),
		Root:      viper.GetString(cfg.KeyReposRoot),
		URL:       url,
	}
	return pkg.Get(config)
}

func runGet(args []string) {
	// Initialize configuration
	cfg.Init(&git.ConfigGlobal{})

	// Create and execute the get command
	cmd := newGetCommand()

	// Set args for cobra to parse
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
