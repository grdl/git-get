package main

import (
	"fmt"
	"git-get/pkg"
	"os"

	"github.com/spf13/cobra"
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
	Args:    cobra.ExactArgs(1),
	Version: fmt.Sprintf("%s - %s, build at %s", version, commit, date),
}

func init() {
	pkg.LoadConf()
}

func Run(cmd *cobra.Command, args []string) {
	url, err := pkg.ParseURL(args[0])
	exitIfError(err)

	_, err = pkg.CloneRepo(url, pkg.Cfg.ReposRoot(), false)
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
