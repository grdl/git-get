// Package cfg provides common configuration to all commands.
// It contains config key names, default values and provides methods to read values from global gitconfig file.
package cfg

import (
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

// ConfigDefaults is a struct with default values for config keys.
type ConfigDefaults struct {
	DefaultHost   string
	Output        string
	ReposRoot     []string
	DefaultScheme string
}

// Defaults is an instance of ConfigDefaults with default values.
var Defaults = ConfigDefaults{
	DefaultHost:   "github.com",
	Output:        OutTree,
	ReposRoot:     []string{fmt.Sprintf("~%c%s", filepath.Separator, "repositories")},
	DefaultScheme: "ssh",
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

	return "git-get " + version
}

// Gitconfig represents gitconfig file.
type Gitconfig interface {
	Get(key string) string
	GetAll(key string) []string
}

// Init initializes viper config registry. Values are looked up in the following order: cli flag, env variable, gitconfig file, default value.
func Init(cfg Gitconfig) {
	readGitconfig(cfg)

	viper.SetEnvPrefix(strings.ToUpper(GitgetPrefix))
	viper.AutomaticEnv()

	// Handle GITGET_ROOT as a path list (like PATH), using OS separator.
	envKey := fmt.Sprintf("%s_%s", strings.ToUpper(GitgetPrefix), strings.ToUpper(KeyReposRoot))
	if val := os.Getenv(envKey); val != "" {
		viper.Set(KeyReposRoot, filepath.SplitList(val))
	}
}

// readGitConfig loads values from gitconfig file into viper's registry.
// Viper doesn't support the gitconfig format so we load it using "git config --global" command and populate a temporary "env" string,
// which is then feed to Viper.
func readGitconfig(cfg Gitconfig) {
	// Root is a list of roots, so it needs to be handled separately using GetAll.
	if val := cfg.GetAll(fmt.Sprintf("%s.%s", GitgetPrefix, KeyReposRoot)); len(val) > 0 {
		viper.SetDefault(KeyReposRoot, val)
	}

	// For other keys we use Get.
	keys := []string{KeyDefaultHost, KeyOutput, KeyDefaultScheme}
	for _, key := range keys {
		if val := cfg.Get(fmt.Sprintf("%s.%s", GitgetPrefix, key)); val != "" {
			viper.SetDefault(key, val)
		}
	}

	// TODO: A hacky way to read boolean flag from gitconfig. Find a cleaner way.
	if val := cfg.Get(fmt.Sprintf("%s.%s", GitgetPrefix, KeySkipHost)); strings.ToLower(val) == "true" {
		viper.SetDefault(KeySkipHost, true)
	}
}

// Expand applies the variables expansion to a viper config of given key.
// If expansion fails or is not needed, the config is not modified.
func Expand(key string) {
	if key == KeyReposRoot {
		roots := viper.GetStringSlice(KeyReposRoot)

		expandedRoots := make([]string, 0, len(roots))
		for _, path := range roots {
			if after, ok := strings.CutPrefix(path, "~"); ok {
				if homeDir, err := os.UserHomeDir(); err == nil {
					path = filepath.Join(homeDir, after)
				}
			}

			expandedRoots = append(expandedRoots, path)
		}

		viper.Set(key, expandedRoots)

		return
	}

	path := viper.GetString(key)
	if after, ok := strings.CutPrefix(path, "~"); ok {
		if homeDir, err := os.UserHomeDir(); err == nil {
			expanded := filepath.Join(homeDir, after)
			viper.Set(key, expanded)
		}
	}
}
