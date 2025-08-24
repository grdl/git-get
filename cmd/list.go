package main

import (
	"fmt"
	"git-get/pkg"
	"git-get/pkg/cfg"
	"git-get/pkg/git"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "git list",
		Short:        "List all repositories cloned by 'git get' and their status.",
		RunE:         runListCommand,
		Args:         cobra.NoArgs,
		Version:      cfg.Version(),
		SilenceUsage: true, // We don't want to show usage on legit errors (eg, wrong path, repo already existing etc.)
	}

	cmd.PersistentFlags().BoolP(cfg.KeyFetch, "f", false, "First fetch from remotes before listing repositories.")
	cmd.PersistentFlags().StringP(cfg.KeyOutput, "o", cfg.Defaults[cfg.KeyOutput], fmt.Sprintf("Output format. Allowed values: [%s].", strings.Join(cfg.AllowedOut, ", ")))
	cmd.PersistentFlags().StringP(cfg.KeyReposRoot, "r", cfg.Defaults[cfg.KeyReposRoot], "Path to repos root where repositories are cloned.")
	cmd.PersistentFlags().BoolP("help", "h", false, "Print this help and exit.")
	cmd.PersistentFlags().BoolP("version", "v", false, "Print version and exit.")

	viper.BindPFlag(cfg.KeyFetch, cmd.PersistentFlags().Lookup(cfg.KeyFetch))
	viper.BindPFlag(cfg.KeyOutput, cmd.PersistentFlags().Lookup(cfg.KeyOutput))
	viper.BindPFlag(cfg.KeyReposRoot, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))

	return cmd
}

func runListCommand(_ *cobra.Command, _ []string) error {
	cfg.Expand(cfg.KeyReposRoot)

	config := &pkg.ListCfg{
		Fetch:  viper.GetBool(cfg.KeyFetch),
		Output: viper.GetString(cfg.KeyOutput),
		Root:   viper.GetString(cfg.KeyReposRoot),
	}

	return pkg.List(config)
}

func runList(args []string) {
	// Initialize configuration
	cfg.Init(&git.ConfigGlobal{})

	// Create and execute the list command
	cmd := newListCommand()

	// Set args for cobra to parse
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
