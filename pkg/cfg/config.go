// Package cfg provides common configuration to all commands.
// It contains config key names, default values and provides methods to read values from global gitconfig file.
package cfg

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// GitgetPrefix is the name of the gitconfig section name and the env var prefix.
const GitgetPrefix = "gitget"

// CLI flag keys.
var (
	KeyBranch      = "branch"
	KeyDump        = "dump"
	KeyDefaultHost = "host"
	KeyFetch       = "fetch"
	KeyOutput      = "out"
	KeyReposRoot   = "root"
)

// Defaults is a map of default values for config keys.
var Defaults = map[string]string{
	KeyDefaultHost: "github.com",
	KeyOutput:      OutTree,
	KeyReposRoot:   fmt.Sprintf("~%c%s", filepath.Separator, "repositories"),
}

// Values for the --out flag.
const (
	OutDump = "dump"
	OutFlat = "flat"
	OutTree = "tree"
)

// AllowedOut are allowed values for the --out flag.
var AllowedOut = []string{OutDump, OutFlat, OutTree}

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
	expandValues()
}

// setMissingValues checks if config values are provided by flags or env vars. If not, it tries loading them from gitconfig file.
// If that fails, the default values are used.
func setMissingValues(cfg Gitconfig) {
	for key, def := range Defaults {
		if isUnsetOrEmpty(key) {
			viper.Set(key, getOrDef(cfg, key, def))
		}
	}
}

func isUnsetOrEmpty(key string) bool {
	return !viper.IsSet(key) || strings.TrimSpace(viper.GetString(key)) == ""
}

func getOrDef(cfg Gitconfig, key string, def string) string {
	if val := cfg.Get(fmt.Sprintf("%s.%s", GitgetPrefix, key)); val != "" {
		return val
	}
	return def
}

// expandValues applies the homedir expansion to a config value. If expansion is not needed value is not modified.
func expandValues() {
	for _, key := range viper.AllKeys() {
		if expanded, err := homedir.Expand(viper.GetString(key)); err == nil {
			viper.Set(key, expanded)
		}
	}
}
