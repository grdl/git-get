// Package cfg provides common configuration to all commands.
// It contains config key names, default values and provides methods to read values from global gitconfig file.
package cfg

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// GitgetPrefix is the name of the gitconfig section name and the env var prefix.
const GitgetPrefix = "gitget"

// CLI flag keys.
var (
	KeyBranch        = "branch"
	KeyDump          = "dump"
	KeyDefaultHost   = "host"
	KeyFetch         = "fetch"
	KeyOutput        = "out"
	KeyDefaultScheme = "scheme"
	KeySkipHost      = "skip-host"
	KeyReposRoot     = "root"
)

// Defaults is a map of default values for config keys.
var Defaults = map[string]string{
	KeyDefaultHost:   "github.com",
	KeyOutput:        OutTree,
	KeyReposRoot:     fmt.Sprintf("~%c%s", filepath.Separator, "repositories"),
	KeyDefaultScheme: "ssh",
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
)

// Version returns a string with version metadata: version number and git commit.
// It returns "git-get development" if version variables are not set during the build.
func Version() string {
	if version == "" {
		return "git-get development"
	}

	if commit != "" {
		return fmt.Sprintf("git-get %s (%s)", version, commit[:7])
	}

	return fmt.Sprintf("git-get %s", version)
}

// Gitconfig represents gitconfig file.
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

	// TODO: Can we somehow iterate over all possible flags?
	for key := range Defaults {
		if val := cfg.Get(fmt.Sprintf("%s.%s", GitgetPrefix, key)); val != "" {
			lines = append(lines, fmt.Sprintf("%s=%s", key, val))
		}
	}

	viper.SetConfigType("env")

	if err := viper.ReadConfig(bytes.NewBuffer([]byte(strings.Join(lines, "\n")))); err != nil {
		// Log error but don't fail - configuration is optional
		fmt.Fprintf(os.Stderr, "Warning: failed to read git config: %v\n", err)
	}

	// TODO: A hacky way to read boolean flag from gitconfig. Find a cleaner way.
	if val := cfg.Get(fmt.Sprintf("%s.%s", GitgetPrefix, KeySkipHost)); strings.ToLower(val) == "true" {
		viper.Set(KeySkipHost, true)
	}
}

// Expand applies the variables expansion to a viper config of given key.
// If expansion fails or is not needed, the config is not modified.
func Expand(key string) {
	path := viper.GetString(key)
	if strings.HasPrefix(path, "~") {
		if homeDir, err := os.UserHomeDir(); err == nil {
			expanded := filepath.Join(homeDir, strings.TrimPrefix(path, "~"))
			viper.Set(key, expanded)
		}
	}
}
