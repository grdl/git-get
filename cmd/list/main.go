package main

import (
	"fmt"
	"git-get/pkg/cfg"
	"git-get/pkg/path"
	"git-get/pkg/print"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmd = &cobra.Command{
	Use:     "git-list",
	Short:   "git list",
	Run:     Run,
	Args:    cobra.NoArgs,
	Version: cfg.Version(),
}

func init() {
	cmd.PersistentFlags().StringP(cfg.KeyReposRoot, "r", "", "repos root")
	cmd.PersistentFlags().StringP(cfg.KeyPrivateKey, "p", "", "SSH private key path")

	cmd.PersistentFlags().BoolP(cfg.KeyFetch, "f", false, "Fetch from remotes when listing repositories")
	cmd.PersistentFlags().StringP(cfg.KeyOutput, "o", cfg.DefOutput, "output format.")

	viper.BindPFlag(cfg.KeyReposRoot, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))
	viper.BindPFlag(cfg.KeyPrivateKey, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))
	viper.BindPFlag(cfg.KeyFetch, cmd.PersistentFlags().Lookup(cfg.KeyFetch))
	viper.BindPFlag(cfg.KeyOutput, cmd.PersistentFlags().Lookup(cfg.KeyOutput))
}

func Run(cmd *cobra.Command, args []string) {
	cfg.InitConfig()

	root := viper.GetString(cfg.KeyReposRoot)

	// TODO: move it to OpenAll and don't export
	paths, err := path.FindRepos()
	exitIfError(err)

	repos, err := path.OpenAll(paths)
	exitIfError(err)

	var printer print.Printer
	switch viper.GetString(cfg.KeyOutput) {
	case cfg.OutFlat:
		printer = &print.FlatPrinter{}
	case cfg.OutTree:
		printer = &print.SimpleTreePrinter{}
	case cfg.OutSmart:
		printer = &print.SmartTreePrinter{}
	case cfg.OutDump:
		printer = &print.DumpPrinter{}
	default:
		err = fmt.Errorf("invalid --out flag; allowed values: %v", []string{cfg.OutFlat, cfg.OutTree, cfg.OutSmart})
	}
	exitIfError(err)

	fmt.Println(printer.Print(root, repos))
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func exitIfError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
