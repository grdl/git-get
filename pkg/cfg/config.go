// Package cfg provides common configuration to all commands.
// It contains config key names, default values and provides methods to read values from global gitconfig file.
package cfg

import (
	"fmt"
	"path"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// GitgetPrefix is the name of the gitconfig section name and the env var prefix.
const GitgetPrefix = "gitget"

// CLI flag keys and their default values.
const (
	KeyBranch = "branch"
	// DefBranch      = "master"
	KeyDump        = "dump"
	KeyDefaultHost = "host"
	DefDefaultHost = "github.com"
	KeyFetch       = "fetch"
	KeyOutput      = "out"
	DefOutput      = OutTree
	KeyPrivateKey  = "privateKey"
	DefPrivateKey  = "id_rsa"
	KeyReposRoot   = "root"
	DefReposRoot   = "repositories"
)

// Values for the --out flag.
const (
	OutDump  = "dump"
	OutFlat  = "flat"
	OutSmart = "smart"
	OutTree  = "tree"
)

// AllowedOut are allowed values for the --out flag.
var AllowedOut = []string{OutDump, OutFlat, OutSmart, OutTree}

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

// Gitconfig represents gitconfig file
type Gitconfig interface {
	Get(key string) string
}

// Init initializes viper config registry. Values are looked up in the following order: cli flag, env variable, gitconfig file, default value.
// Viper doesn't support gitconfig file format so it can't find missing values there automatically. They need to be specified in setMissingValues func.
//
// Because it reads the cli flags it needs to be called after the cmd.Execute().
func Init(cfg Gitconfig) {
	viper.SetEnvPrefix(strings.ToUpper(GitgetPrefix))
	viper.AutomaticEnv()

	setMissingValues(cfg)
}

// setMissingValues checks if config values are provided by flags or env vars. If not, it tries loading them from gitconfig file.
// If that fails, the default values are used.
func setMissingValues(cfg Gitconfig) {
	if isUnsetOrEmpty(KeyReposRoot) {
		viper.Set(KeyReposRoot, getOrDef(cfg, KeyReposRoot, path.Join(home(), DefReposRoot)))
	}

	if isUnsetOrEmpty(KeyDefaultHost) {
		viper.Set(KeyDefaultHost, getOrDef(cfg, KeyDefaultHost, DefDefaultHost))
	}

	if isUnsetOrEmpty(KeyPrivateKey) {
		viper.Set(KeyPrivateKey, getOrDef(cfg, KeyPrivateKey, path.Join(home(), ".ssh", DefPrivateKey)))
	}
}

func getOrDef(cfg Gitconfig, key string, def string) string {
	if val := cfg.Get(key); val != "" {
		return val
	}
	return def
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
