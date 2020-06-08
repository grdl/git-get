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
	Args:    cobra.MaximumNArgs(1), // TODO: add custom validator
	Version: fmt.Sprintf("%s - %s, build at %s", version, commit, date),
}

func init() {
	cmd.PersistentFlags().BoolP(cfg.KeyList, "l", false, "Lists all repositories inside git-get root")
	cmd.PersistentFlags().BoolP(cfg.KeyFetch, "f", false, "Fetch from remotes when listing repositories")
	cmd.PersistentFlags().StringP(cfg.KeyReposRoot, "r", "", "repos root")
	cmd.PersistentFlags().StringP(cfg.KeyPrivateKey, "p", "", "SSH private key path")
	cmd.PersistentFlags().StringP(cfg.KeyOutput, "o", cfg.DefOutput, "output format.")
	cmd.PersistentFlags().StringP(cfg.KeyBranch, "b", cfg.DefBranch, "Branch (or tag) to checkout after cloning")

	viper.BindPFlag(cfg.KeyList, cmd.PersistentFlags().Lookup(cfg.KeyList))
	viper.BindPFlag(cfg.KeyFetch, cmd.PersistentFlags().Lookup(cfg.KeyFetch))
	viper.BindPFlag(cfg.KeyReposRoot, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))
	viper.BindPFlag(cfg.KeyPrivateKey, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))
	viper.BindPFlag(cfg.KeyOutput, cmd.PersistentFlags().Lookup(cfg.KeyOutput))
	viper.BindPFlag(cfg.KeyBranch, cmd.PersistentFlags().Lookup(cfg.KeyBranch))
}

func Run(cmd *cobra.Command, args []string) {
	cfg.InitConfig()

	root := viper.GetString(cfg.KeyReposRoot)
	if viper.GetBool(cfg.KeyList) {
		// TODO: move it to OpenAll and don't export
		paths, err := path.FindRepos()
		exitIfError(err)

		repos, err := path.OpenAll(paths)
		exitIfError(err)

		var printer print.Printer
		switch viper.GetString(cfg.KeyOutput) {
		case cfg.OutFlat:
			printer = &print.FlatPrinter{}
		case cfg.OutSimple:
			printer = &print.SimpleTreePrinter{}
		case cfg.OutSmart:
			printer = &print.SmartTreePrinter{}
		default:
			err = fmt.Errorf("invalid --out flag; allowed values: %v", []string{cfg.OutFlat, cfg.OutSimple, cfg.OutSmart})
		}
		exitIfError(err)

		fmt.Println(printer.Print(root, repos))

		os.Exit(0)
	}

	url, err := path.ParseURL(args[0])
	exitIfError(err)

	branch := viper.GetString(cfg.KeyBranch)
	repoPath := pathpkg.Join(root, path.URLToPath(url))
	_, err = git.CloneRepo(url, repoPath, branch, false)
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
