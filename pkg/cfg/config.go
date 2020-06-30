// Package cfg provides common configuration to all commands.
// It contains config key names, default values and provides methods to read values from global gitconfig file.
package cfg

import (
	"bytes"
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
func Init(cfg Gitconfig) {
	readGitconfig(cfg)

	viper.SetEnvPrefix(strings.ToUpper(GitgetPrefix))
	viper.AutomaticEnv()
}

// readGitConfig loads values from gitconfig file into viper's registry.
// Viper doesn't support the gitconfig format so we load it using "git config --global" command and populate a temporary "env" string,
// which is then feed to Viper.
func readGitconfig(cfg Gitconfig) {
	var lines []string

	for key := range Defaults {
		if val := cfg.Get(fmt.Sprintf("%s.%s", GitgetPrefix, key)); val != "" {
			lines = append(lines, fmt.Sprintf("%s=%s", key, val))
		}
	}

	viper.SetConfigType("env")
	viper.ReadConfig(bytes.NewBuffer([]byte(strings.Join(lines, "\n"))))
}

// Expand applies the variables expansion to a viper config of given key.
// If expansion fails or is not needed, the config is not modified.
func Expand(key string) {
	if expanded, err := homedir.Expand(viper.GetString(key)); err == nil {
		viper.Set(key, expanded)
	}
}
