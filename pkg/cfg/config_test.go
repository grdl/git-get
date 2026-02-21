package cfg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	envVarName    = strings.ToUpper(fmt.Sprintf("%s_%s", GitgetPrefix, KeyDefaultHost))
	fromGitconfig = "value.from.gitconfig"
	fromEnv       = "value.from.env"
	fromFlag      = "value.from.flag"
)

//nolint:paralleltest // These tests modify global state (viper, env vars) and cannot run in parallel
func TestConfig(t *testing.T) {
	tests := []struct {
		name        string
		configMaker func(*testing.T)
		key         string
		want        string
	}{
		{
			name:        "no config",
			configMaker: testConfigEmpty,
			key:         KeyDefaultHost,
			want:        Defaults.DefaultHost,
		},
		{
			name:        "value only in gitconfig",
			configMaker: testConfigOnlyInGitconfig,
			key:         KeyDefaultHost,
			want:        fromGitconfig,
		},
		{
			name:        "value only in env var",
			configMaker: testConfigOnlyInEnvVar,
			key:         KeyDefaultHost,
			want:        fromEnv,
		},
		{
			name:        "value in gitconfig and env var",
			configMaker: testConfigInGitconfigAndEnvVar,
			key:         KeyDefaultHost,
			want:        fromEnv,
		},
		{
			name:        "value in flag",
			configMaker: testConfigInFlag,
			key:         KeyDefaultHost,
			want:        fromFlag,
		},
		{
			name:        "multiple roots in gitconfig",
			configMaker: testConfigMultipleRootsInGitconfig,
			key:         KeyReposRoot,
			want:        "root1,root2",
		},
		{
			name:        "multiple roots in env var",
			configMaker: testConfigMultipleRootsInEnvVar,
			key:         KeyReposRoot,
			want:        "root1,root2",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set manual defaults for test keys
			switch test.key {
			case KeyDefaultHost:
				viper.SetDefault(test.key, Defaults.DefaultHost)
			case KeyReposRoot:
				viper.SetDefault(test.key, Defaults.ReposRoot)
			}

			test.configMaker(t)

			var got string
			if test.key == KeyReposRoot {
				got = strings.Join(viper.GetStringSlice(test.key), ",")
			} else {
				got = viper.GetString(test.key)
			}

			if got != test.want {
				t.Errorf("expected %q; got %q", test.want, got)
			}

			// Clear env variables and reset viper registry after each test so they impact other tests.
			os.Clearenv()
			viper.Reset()
		})
	}
}

type gitconfigEmpty struct{}

func (c *gitconfigEmpty) Get(key string) string {
	return ""
}

func (c *gitconfigEmpty) GetAll(key string) []string {
	return nil
}

type gitconfigValid struct{}

func (c *gitconfigValid) Get(key string) string {
	return fromGitconfig
}

func (c *gitconfigValid) GetAll(key string) []string {
	return []string{fromGitconfig}
}

type gitconfigMultipleRoots struct{}

func (c *gitconfigMultipleRoots) Get(key string) string {
	return "root1"
}

func (c *gitconfigMultipleRoots) GetAll(key string) []string {
	return []string{"root1", "root2"}
}

func testConfigEmpty(t *testing.T) {
	t.Helper()
	Init(&gitconfigEmpty{})
}

func testConfigOnlyInGitconfig(t *testing.T) {
	t.Helper()
	Init(&gitconfigValid{})
}

func testConfigOnlyInEnvVar(t *testing.T) {
	t.Helper()
	Init(&gitconfigEmpty{})
	t.Setenv(envVarName, fromEnv)
}

func testConfigInGitconfigAndEnvVar(t *testing.T) {
	t.Helper()
	Init(&gitconfigValid{})
	t.Setenv(envVarName, fromEnv)
}

func testConfigInFlag(t *testing.T) {
	t.Helper()
	Init(&gitconfigValid{})
	t.Setenv(envVarName, fromEnv)

	cmd := cobra.Command{}
	cmd.PersistentFlags().String(KeyDefaultHost, Defaults.DefaultHost, "")

	if err := viper.BindPFlag(KeyDefaultHost, cmd.PersistentFlags().Lookup(KeyDefaultHost)); err != nil {
		t.Fatalf("failed to bind flag: %v", err)
	}

	cmd.SetArgs([]string{"--" + KeyDefaultHost, fromFlag})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to execute command: %v", err)
	}
}

func testConfigMultipleRootsInGitconfig(t *testing.T) {
	t.Helper()
	Init(&gitconfigMultipleRoots{})
}

func testConfigMultipleRootsInEnvVar(t *testing.T) {
	t.Helper()

	envKey := fmt.Sprintf("%s_%s", strings.ToUpper(GitgetPrefix), strings.ToUpper(KeyReposRoot))
	t.Setenv(envKey, strings.Join([]string{"root1", "root2"}, string(filepath.ListSeparator)))
	Init(&gitconfigEmpty{})
}
