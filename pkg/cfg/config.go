package cfg

import (
	"fmt"
	"path"
	"strings"

	"github.com/go-git/go-git/v5/config"
	plumbing "github.com/go-git/go-git/v5/plumbing/format/config"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Gitconfig section name and env var prefix
const GitgetPrefix = "gitget"

// Flag keys and their default values
const (
	KeyBranch      = "branch"
	DefBranch      = "master"
	KeyBundle      = "bundle"
	KeyDefaultHost = "defaultHost"
	DefDefaultHost = "github.com"
	KeyFetch       = "fetch"
	KeyList        = "list"
	KeyOutput      = "out"
	DefOutput      = OutFlat
	KeyPrivateKey  = "privateKey"
	DefPrivateKey  = "id_rsa"
	KeyReposRoot   = "reposRoot"
	DefReposRoot   = "repositories"
)

// Allowed values for the --out flag
const (
	OutFlat   = "flat"
	OutSmart  = "smart"
	OutSimple = "simple"
)

// Version metadata set by ldflags during the build.
var (
	version string
	commit  string
	date    string
)

// Version returns a string with version metadata: version number, git sha and build date.
// It returns "development" if version variables are not set during the build.
func Version() string {
	if version == "" {
		return "development"
	}

	return fmt.Sprintf("%s - revision %s built at %s", version, commit[:6], date)
}

// gitconfig provides methods for looking up configiration values inside .gitconfig file
type gitconfig struct {
	*config.Config
}

// InitConfig initializes viper config registry. Values are looked up in the following order: cli flag, env variable, gitconfig file, default value
// Viper doesn't support gitconfig file format so it can't find missing values there automatically. They need to be specified in setMissingValues func.
//
// Because it reads the cli flags it needs to be called after the cmd.Execute().
func InitConfig() {
	viper.SetEnvPrefix(strings.ToUpper(GitgetPrefix))
	viper.AutomaticEnv()

	cfg := loadGitconfig()
	setMissingValues(cfg)
}

// loadGitconfig loads configuration from a gitconfig file.
// We ignore errors when gitconfig file can't be found, opened or parsed. In those cases viper will provide default config values.
func loadGitconfig() *gitconfig {
	// TODO: load system scope
	cfg, _ := config.LoadConfig(config.GlobalScope)

	return &gitconfig{
		Config: cfg,
	}
}

// setMissingValues checks if config values are provided by flags or env vars. If not, it tries loading them from gitconfig file.
// If that fails, the default values are used.
func setMissingValues(cfg *gitconfig) {
	if isUnsetOrEmpty(KeyReposRoot) {
		viper.Set(KeyReposRoot, cfg.get(KeyReposRoot, path.Join(home(), DefReposRoot)))
	}

	if isUnsetOrEmpty(KeyDefaultHost) {
		viper.Set(KeyDefaultHost, cfg.get(KeyDefaultHost, DefDefaultHost))
	}

	if isUnsetOrEmpty(KeyPrivateKey) {
		viper.Set(KeyPrivateKey, cfg.get(KeyPrivateKey, path.Join(home(), ".ssh", DefPrivateKey)))
	}
}

// get looks up the value for a given key in gitconfig file.
// It returns the default value when gitconfig is missing, or it doesn't contain a gitget section,
// or if the section is empty, or if it doesn't contain a valid value for the key.
func (c *gitconfig) get(key string, def string) string {
	if c == nil || c.Config == nil {
		return def
	}

	gitget := c.findGitconfigSection(GitgetPrefix)
	if gitget == nil {
		return def
	}

	opt := gitget.Option(key)
	if strings.TrimSpace(opt) == "" {
		return def
	}

	return opt
}

func (c *gitconfig) findGitconfigSection(name string) *plumbing.Section {
	for _, s := range c.Raw.Sections {
		if strings.ToLower(s.Name) == strings.ToLower(name) {
			return s
		}
	}

	return nil
}

// home returns path to a home directory or empty string if can't be found.
// Using empty string means that in the unlikely situation where home dir can't be found
// and there's no reposRoot specified by any of the config methods, the current dir will be used as repos root.
func home() string {
	home, err := homedir.Dir()
	if err != nil {
		return ""
	}

	return home
}

func isUnsetOrEmpty(key string) bool {
	return !viper.IsSet(key) || strings.TrimSpace(viper.GetString(key)) == ""
}
