package main

import (
	"fmt"
	"git-get/pkg"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

var cmd = &cobra.Command{
	Use:     "git-get <repo>",
	Short:   "git get",
	Run:     Run,
	Args:    cobra.MaximumNArgs(1),
	Version: fmt.Sprintf("%s - %s, build at %s", version, commit, date),
}

var list bool

func init() {
	cmd.PersistentFlags().BoolVarP(&list, "list", "l", false, "Lists all repositories inside git-get root")
	cmd.PersistentFlags().StringP(pkg.KeyReposRoot, "r", "", "repos root")
	cmd.PersistentFlags().StringP(pkg.KeyPrivateKey, "p", "", "SSH private key path")
	viper.BindPFlag(pkg.KeyReposRoot, cmd.PersistentFlags().Lookup(pkg.KeyReposRoot))
	viper.BindPFlag(pkg.KeyPrivateKey, cmd.PersistentFlags().Lookup(pkg.KeyReposRoot))
}

func Run(cmd *cobra.Command, args []string) {
	pkg.InitConfig()

	if list {
		paths, err := pkg.FindRepos()
		exitIfError(err)

		repos, err := pkg.OpenAll(paths)
		exitIfError(err)

		pkg.PrintRepos(repos)
		os.Exit(0)
	}

	url, err := pkg.ParseURL(args[0])
	exitIfError(err)

	_, err = pkg.CloneRepo(url, viper.GetString(pkg.KeyReposRoot), false)
	exitIfError(err)
}

func main() {
	err := cmd.Execute()
	exitIfError(err)
}

func exitIfError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
