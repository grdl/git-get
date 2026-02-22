package main

import (
	"fmt"
	"os"

	"github.com/grdl/git-get/pkg/cfg"
	"github.com/grdl/git-get/pkg/git"

	"github.com/spf13/cobra"
)

const initLong = `Print the shell initialization code to enable 'git cd' and 'gitcd' functionality.
Since a child process cannot change the shell's directory, a shell function wrapper is required.

Add the following to your shell profile (e.g., ~/.bashrc, ~/.zshrc, or ~/.config/fish/config.fish):

Bash/Zsh:
  eval "$(git-get shell-init bash)"

Fish:
  git-get shell-init fish | source
`

func newInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "shell-init [bash|zsh|fish]",
		Short: "Print shell initialization code for 'git cd'.",
		Long:  initLong,
		Args:  cobra.MaximumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			shell := "bash"
			if len(args) > 0 {
				shell = args[0]
			}

			switch shell {
			case "bash", "zsh":
				fmt.Print(`
# git-get shell integration
gitcd() {
    local target
    target=$(git-get cd "$@")
    if [ -n "$target" ]; then
        cd "$target"
    fi
}

# Optional: Override git to support 'git cd'
git() {
    if [ "$1" = "cd" ]; then
        shift
        gitcd "$@"
    else
        command git "$@"
    fi
}
`)
			case "fish":
				//nolint:dupword
				fmt.Print(`
# git-get shell integration
function gitcd
    set -l target (git-get cd $argv)
    if test -n "$target"
        cd $target
    end
end

# Optional: Override git to support 'git cd'
function git
    if test "$argv[1]" = "cd"
        set -e argv[1]
        gitcd $argv
    else
        command git $argv
    end
end
`)
			default:
				fmt.Fprintf(os.Stderr, "Unsupported shell: %s\n", shell)
				os.Exit(1)
			}
		},
	}
}

func runInit(args []string) {
	cfg.Init(&git.ConfigGlobal{})

	cmd := newInitCommand()

	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
