package main

import (
	"fmt"
	"os"

	"github.com/grdl/git-get/pkg"
	"github.com/grdl/git-get/pkg/cfg"
	"github.com/grdl/git-get/pkg/git"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cdLong = `Fuzzy find a repository and print its absolute path.
Since a child process cannot change the current working directory of the shell,
this command prints the path to stdout.

To use it effectively, add a shell function or alias to your shell profile:

Bash/Zsh:
  gcd() {
    local target
    target=$(git-get cd "$@")
    if [ -n "$target" ]; then
      cd "$target"
    fi
  }

Fish:
  function gcd
    set target (git-get cd $argv)
    if test -n "$target"
      cd $target
    end
  end
`

func newCdCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "cd [QUERY]",
		Short:        "Fuzzy find a repository and print its absolute path.",
		Long:         cdLong,
		RunE:         runCdCommand,
		Args:         cobra.MaximumNArgs(1),
		Version:      cfg.Version(),
		SilenceUsage: true,
	}

	cmd.PersistentFlags().StringSliceP(cfg.KeyReposRoot, "r", cfg.Defaults.ReposRoot, "Path to repos root where repositories are scanned.")

	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		panic(fmt.Sprintf("failed to bind flags: %v", err))
	}

	return cmd
}

func runCdCommand(_ *cobra.Command, args []string) error {
	cfg.Expand(cfg.KeyReposRoot)

	var query string
	if len(args) > 0 {
		query = args[0]
	}

	config := &pkg.CdCfg{
		Roots: viper.GetStringSlice(cfg.KeyReposRoot),
		Query: query,
	}

	return pkg.Cd(config)
}

func runCd(args []string) {
	// Initialize configuration
	cfg.Init(&git.ConfigGlobal{})

	// Create and execute the cd command
	cmd := newCdCommand()

	// Set args for cobra to parse
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
