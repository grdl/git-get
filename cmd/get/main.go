package main

import (
	"fmt"
	"git-get/pkg/cfg"
	"git-get/pkg/git"
	"git-get/pkg/path"
	"os"
	pathpkg "path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmd = &cobra.Command{
	Use:     "git-get <repo>",
	Short:   "git get",
	Run:     Run,
	Args:    cobra.MaximumNArgs(1), // TODO: add custom validator
	Version: cfg.Version(),
}

func init() {
	cmd.PersistentFlags().StringP(cfg.KeyReposRoot, "r", "", "repos root")
	cmd.PersistentFlags().StringP(cfg.KeyPrivateKey, "p", "", "SSH private key path")
	cmd.PersistentFlags().StringP(cfg.KeyBundle, "u", "", "Bundle file path")

	cmd.PersistentFlags().StringP(cfg.KeyBranch, "b", cfg.DefBranch, "Branch (or tag) to checkout after cloning")

	viper.BindPFlag(cfg.KeyReposRoot, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))
	viper.BindPFlag(cfg.KeyPrivateKey, cmd.PersistentFlags().Lookup(cfg.KeyReposRoot))
	viper.BindPFlag(cfg.KeyBundle, cmd.PersistentFlags().Lookup(cfg.KeyBundle))
	viper.BindPFlag(cfg.KeyBranch, cmd.PersistentFlags().Lookup(cfg.KeyBranch))
}

func Run(cmd *cobra.Command, args []string) {
	cfg.InitConfig()

	root := viper.GetString(cfg.KeyReposRoot)

	if bundle := viper.GetString(cfg.KeyBundle); bundle != "" {
		opts, err := path.ParseBundleFile(bundle)
		exitIfError(err)

		for _, opt := range opts {
			path := pathpkg.Join(root, path.URLToPath(opt.URL))
			opt.Path = path
			_, _ = git.CloneRepo(opt)
		}
		os.Exit(0)
	}

	url, err := path.ParseURL(args[0])
	exitIfError(err)

	branch := viper.GetString(cfg.KeyBranch)
	path := pathpkg.Join(root, path.URLToPath(url))

	cloneOpts := &git.CloneOpts{
		URL:    url,
		Path:   path,
		Branch: branch,
	}

	_, err = git.CloneRepo(cloneOpts)
	exitIfError(err)
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
