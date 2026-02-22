package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/grdl/git-get/pkg"
	"github.com/grdl/git-get/pkg/cfg"
	"github.com/grdl/git-get/pkg/git"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "git list [QUERY]",
		Short:        "List all repositories cloned by 'git get' and their status.",
		Long:         "List all repositories cloned by 'git get' and their status. If [QUERY] is provided, filter repositories using fuzzy search.",
		RunE:         runListCommand,
		Args:         cobra.MaximumNArgs(1),
		Version:      cfg.Version(),
		SilenceUsage: true, // We don't want to show usage on legit errors (eg, wrong path, repo already existing etc.)
	}

	cmd.PersistentFlags().BoolP(cfg.KeyFetch, "f", false, "First fetch from remotes before listing repositories.")
	cmd.PersistentFlags().StringP(cfg.KeyOutput, "o", cfg.Defaults.Output, fmt.Sprintf("Output format. Allowed values: [%s].", strings.Join(cfg.AllowedOut, ", ")))
	cmd.PersistentFlags().StringSliceP(cfg.KeyReposRoot, "r", cfg.Defaults.ReposRoot, "Path to repos root where repositories are scanned.")
	cmd.PersistentFlags().BoolP("help", "h", false, "Print this help and exit.")
	cmd.PersistentFlags().BoolP("version", "v", false, "Print version and exit.")

	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		panic(fmt.Sprintf("failed to bind flags: %v", err))
	}

	return cmd
}

func runListCommand(_ *cobra.Command, args []string) error {
	cfg.Expand(cfg.KeyReposRoot)

	var query string
	if len(args) > 0 {
		query = args[0]
	}

	config := &pkg.ListCfg{
		Fetch:  viper.GetBool(cfg.KeyFetch),
		Output: viper.GetString(cfg.KeyOutput),
		Roots:  viper.GetStringSlice(cfg.KeyReposRoot),
		Query:  query,
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
