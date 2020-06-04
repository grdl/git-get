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
var reposRoot string

func init() {
	// pkg.LoadConf()

	cmd.PersistentFlags().BoolVarP(&list, "list", "l", false, "Lists all repositories inside git-get root")
	cmd.PersistentFlags().StringVarP(&reposRoot, "reposRoot", "r", "", "repos root")
	viper.BindPFlag("reposRoot", cmd.PersistentFlags().Lookup("reposRoot"))

	pkg.InitConfig()

}

func Run(cmd *cobra.Command, args []string) {
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
