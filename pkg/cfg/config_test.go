package cfg

import (
	"fmt"
	"os"
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
			want:        Defaults[KeyDefaultHost],
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			viper.SetDefault(test.key, Defaults[KeyDefaultHost])
			test.configMaker(t)

			got := viper.GetString(test.key)
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

type gitconfigValid struct{}

func (c *gitconfigValid) Get(key string) string {
	return fromGitconfig
}

func testConfigEmpty(t *testing.T) {
	Init(&gitconfigEmpty{})
}

func testConfigOnlyInGitconfig(t *testing.T) {
	Init(&gitconfigValid{})
}

func testConfigOnlyInEnvVar(t *testing.T) {
	Init(&gitconfigEmpty{})
	t.Setenv(envVarName, fromEnv)
}

func testConfigInGitconfigAndEnvVar(t *testing.T) {
	Init(&gitconfigValid{})
	t.Setenv(envVarName, fromEnv)
}

func testConfigInFlag(t *testing.T) {
	Init(&gitconfigValid{})
	t.Setenv(envVarName, fromEnv)

	cmd := cobra.Command{}
	cmd.PersistentFlags().String(KeyDefaultHost, Defaults[KeyDefaultHost], "")

	if err := viper.BindPFlag(KeyDefaultHost, cmd.PersistentFlags().Lookup(KeyDefaultHost)); err != nil {
		t.Fatalf("failed to bind flag: %v", err)
	}

	cmd.SetArgs([]string{"--" + KeyDefaultHost, fromFlag})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to execute command: %v", err)
	}
}
