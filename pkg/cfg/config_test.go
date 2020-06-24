package cfg

import (
	"fmt"
	"git-get/pkg_old/cfg"
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
			want:        DefDefaultHost,
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
	os.Setenv(envVarName, fromEnv)

	Init(&gitconfigEmpty{})
}

func testConfigInGitconfigAndEnvVar(t *testing.T) {
	os.Setenv(envVarName, fromEnv)

	Init(&gitconfigValid{})
}

func testConfigInFlag(t *testing.T) {
	os.Setenv(envVarName, fromEnv)

	cmd := cobra.Command{}
	cmd.PersistentFlags().String(cfg.KeyDefaultHost, cfg.DefDefaultHost, "")
	viper.BindPFlag(cfg.KeyDefaultHost, cmd.PersistentFlags().Lookup(cfg.KeyDefaultHost))

	cmd.SetArgs([]string{"--" + cfg.KeyDefaultHost, fromFlag})
	cmd.Execute()
	Init(&gitconfigValid{})
}
