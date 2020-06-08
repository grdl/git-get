package main

import (
	"fmt"
	"git-get/cfg"
	"git-get/git"
	"git-get/path"
	"git-get/print"
	"os"

	pathpkg "path"

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
	cmd.PersistentFlags().StringP(cfg.KeyReposRoot, "r", "", "repos root")
	cmd.PersistentFlags().StringP(cfg.KeyPrivateKey, "p", "", "SSH private key path")
	viper.BindPFlag(cfg.KeyReposRoot, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))
	viper.BindPFlag(cfg.KeyPrivateKey, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))
}

func Run(cmd *cobra.Command, args []string) {
	cfg.InitConfig()

	root := viper.GetString(cfg.KeyReposRoot)
	if list {
		paths, err := path.FindRepos()
		exitIfError(err)

		repos, err := path.OpenAll(paths)
		exitIfError(err)

		//tree := BuildTree(root, repos)
		//fmt.Println(RenderSmartTree(tree))

		printer := print.NewFlatPrinter()
		fmt.Println(printer.Print(root, repos))

		os.Exit(0)
	}

	url, err := path.ParseURL(args[0])
	exitIfError(err)
	repoPath := pathpkg.Join(root, path.URLToPath(url))

	_, err = git.CloneRepo(url, repoPath, false)
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
