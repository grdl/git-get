package cfg

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/viper"
)

const (
	EnvDefaultHost = "GITGET_DEFAULTHOST"
	EnvReposRoot   = "GITGET_REPOSROOT"
)

func newConfigWithFullGitconfig() *gitconfig {
	cfg := config.NewConfig()

	gitget := cfg.Raw.Section(GitgetPrefix)
	gitget.AddOption(KeyReposRoot, "file.root")
	gitget.AddOption(KeyDefaultHost, "file.host")

	return &gitconfig{
		Config: cfg,
	}
}

func newConfigWithEmptyGitgetSection() *gitconfig {
	cfg := config.NewConfig()

	_ = cfg.Raw.Section(GitgetPrefix)

	return &gitconfig{
		Config: cfg,
	}
}

func newConfigWithEmptyValues() *gitconfig {
	cfg := config.NewConfig()

	gitget := cfg.Raw.Section(GitgetPrefix)
	gitget.AddOption(KeyReposRoot, "")
	gitget.AddOption(KeyDefaultHost, "   ")

	return &gitconfig{
		Config: cfg,
	}
}

func newConfigWithoutGitgetSection() *gitconfig {
	cfg := config.NewConfig()

	return &gitconfig{
		Config: cfg,
	}
}

func newConfigWithEmptyGitconfig() *gitconfig {
	return &gitconfig{
		Config: nil,
	}
}

func newConfigWithEnvVars() *gitconfig {
	_ = os.Setenv(EnvDefaultHost, "env.host")
	_ = os.Setenv(EnvReposRoot, "env.root")

	return &gitconfig{
		Config: nil,
	}
}

func newConfigWithGitconfigAndEnvVars() *gitconfig {
	cfg := config.NewConfig()

	gitget := cfg.Raw.Section(GitgetPrefix)
	gitget.AddOption(KeyReposRoot, "file.root")
	gitget.AddOption(KeyDefaultHost, "file.host")

	_ = os.Setenv(EnvDefaultHost, "env.host")
	_ = os.Setenv(EnvReposRoot, "env.root")

	return &gitconfig{
		Config: cfg,
	}
}

func newConfigWithEmptySectionAndEnvVars() *gitconfig {
	cfg := config.NewConfig()

	_ = cfg.Raw.Section(GitgetPrefix)

	_ = os.Setenv(EnvDefaultHost, "env.host")
	_ = os.Setenv(EnvReposRoot, "env.root")

	return &gitconfig{
		Config: cfg,
	}
}

func newConfigWithMixed() *gitconfig {
	cfg := config.NewConfig()

	gitget := cfg.Raw.Section(GitgetPrefix)
	gitget.AddOption(KeyReposRoot, "file.root")
	gitget.AddOption(KeyDefaultHost, "file.host")

	_ = os.Setenv(EnvDefaultHost, "env.host")

	return &gitconfig{
		Config: cfg,
	}
}

func TestConfig(t *testing.T) {
	defReposRoot := path.Join(home(), DefReposRoot)

	var tests = []struct {
		makeConfig      func() *gitconfig
		wantReposRoot   string
		wantDefaultHost string
	}{
		{newConfigWithFullGitconfig, "file.root", "file.host"},
		{newConfigWithoutGitgetSection, defReposRoot, DefDefaultHost},
		{newConfigWithEmptyGitconfig, defReposRoot, DefDefaultHost},
		{newConfigWithEnvVars, "env.root", "env.host"},
		{newConfigWithGitconfigAndEnvVars, "env.root", "env.host"},
		{newConfigWithEmptySectionAndEnvVars, "env.root", "env.host"},
		{newConfigWithEmptyGitgetSection, defReposRoot, DefDefaultHost},
		{newConfigWithEmptyValues, defReposRoot, DefDefaultHost},
		{newConfigWithMixed, "file.root", "env.host"},
	}

	for _, test := range tests {
		viper.SetEnvPrefix(strings.ToUpper(GitgetPrefix))
		viper.AutomaticEnv()

		cfg := test.makeConfig()
		setMissingValues(cfg)

		if viper.GetString(KeyDefaultHost) != test.wantDefaultHost {
			t.Errorf("Wrong %s value, got: %s; want: %s", KeyDefaultHost, viper.GetString(KeyDefaultHost), test.wantDefaultHost)
		}

		if viper.GetString(KeyReposRoot) != test.wantReposRoot {
			t.Errorf("Wrong %s value, got: %s; want: %s", KeyReposRoot, viper.GetString(KeyReposRoot), test.wantReposRoot)
		}

		// Unset env variables and reset viper registry after each test
		viper.Reset()
		err := os.Unsetenv(EnvDefaultHost)
		checkFatal(t, err)
		err = os.Unsetenv(EnvReposRoot)
		checkFatal(t, err)
	}
}

func checkFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("%+v", err)
	}
}
