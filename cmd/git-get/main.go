package main

import (
	"fmt"
	"git-get/pkg"
	"os"

	"github.com/spf13/cobra"
)

const ReposRoot = "/tmp/gitget"

var cmd = &cobra.Command{
	Use:     "git-get",
	Short:   "git get",
	RunE:    Run,
	Version: "0.0.0",
}

func init() {
	//cmd.PersistentFlags().
}

func Run(cmd *cobra.Command, args []string) error {
	url, err := pkg.ParseURL(args[0])
	if err != nil {
		return err
	}

	_, err = pkg.CloneRepo(url, ReposRoot, false)
	return err
}

func Execute() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
