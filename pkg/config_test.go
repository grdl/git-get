package pkg

import (
	"os"
	"path"
	"testing"

	"github.com/go-git/go-git/v5/config"
)

func newConfigWithFullGitconfig() *Conf {
	gitconfig := config.NewConfig()

	gitget := gitconfig.Raw.Section(CfgSection)
	gitget.AddOption(CfgReposRoot, "file.root")
	gitget.AddOption(CfgDefaultHost, "file.host")

	return &Conf{
		gitconfig: gitconfig,
	}
}

func newConfigWithEmptyGitgetSection() *Conf {
	gitconfig := config.NewConfig()

	_ = gitconfig.Raw.Section(CfgSection)

	return &Conf{
		gitconfig: gitconfig,
	}
}

func newConfigWithEmptyValues() *Conf {
	gitconfig := config.NewConfig()

	gitget := gitconfig.Raw.Section(CfgSection)
	gitget.AddOption(CfgReposRoot, "")
	gitget.AddOption(CfgDefaultHost, "   ")

	return &Conf{
		gitconfig: gitconfig,
	}
}

func newConfigWithoutGitgetSection() *Conf {
	gitconfig := config.NewConfig()

	return &Conf{
		gitconfig: gitconfig,
	}
}

func newConfigWithEmptyGitconfig() *Conf {
	return &Conf{
		gitconfig: nil,
	}
}

func newConfigWithEnvVars() *Conf {
	_ = os.Setenv(EnvDefaultHost, "env.host")
	_ = os.Setenv(EnvReposRoot, "env.root")

	return &Conf{
		gitconfig: nil,
	}
}

func newConfigWithGitconfigAndEnvVars() *Conf {
	gitconfig := config.NewConfig()

	gitget := gitconfig.Raw.Section(CfgSection)
	gitget.AddOption(CfgReposRoot, "file.root")
	gitget.AddOption(CfgDefaultHost, "file.host")

	_ = os.Setenv(EnvDefaultHost, "env.host")
	_ = os.Setenv(EnvReposRoot, "env.root")

	return &Conf{
		gitconfig: gitconfig,
	}
}

func newConfigWithEmptySectionAndEnvVars() *Conf {
	gitconfig := config.NewConfig()

	_ = gitconfig.Raw.Section(CfgSection)

	_ = os.Setenv(EnvDefaultHost, "env.host")
	_ = os.Setenv(EnvReposRoot, "env.root")

	return &Conf{
		gitconfig: gitconfig,
	}
}

func newConfigWithMixed() *Conf {
	gitconfig := config.NewConfig()

	gitget := gitconfig.Raw.Section(CfgSection)
	gitget.AddOption(CfgReposRoot, "file.root")
	gitget.AddOption(CfgDefaultHost, "file.host")

	_ = os.Setenv(EnvDefaultHost, "env.host")

	return &Conf{
		gitconfig: gitconfig,
	}
}

func TestConfig(t *testing.T) {
	defReposRoot := path.Join(home(), DefaultReposRootSubpath)

	var tests = []struct {
		makeConfig      func() *Conf
		wantReposRoot   string
		wantDefaultHost string
	}{
		{newConfigWithFullGitconfig, "file.root", "file.host"},
		{newConfigWithoutGitgetSection, defReposRoot, DefaultHost},
		{newConfigWithEmptyGitconfig, defReposRoot, DefaultHost},
		{newConfigWithEnvVars, "env.root", "env.host"},
		{newConfigWithGitconfigAndEnvVars, "env.root", "env.host"},
		{newConfigWithEmptySectionAndEnvVars, "env.root", "env.host"},
		{newConfigWithEmptyGitgetSection, defReposRoot, DefaultHost},
		{newConfigWithEmptyValues, defReposRoot, DefaultHost},
		{newConfigWithMixed, "file.root", "env.host"},
	}

	for _, test := range tests {
		cfg := test.makeConfig()

		if cfg.ReposRoot() != test.wantReposRoot {
			t.Errorf("Wrong reposRoot value, got: %+v; want: %+v", cfg.ReposRoot(), test.wantReposRoot)
		}

		if cfg.DefaultHost() != test.wantDefaultHost {
			t.Errorf("Wrong defaultHost value, got: %+v; want: %+v", cfg.DefaultHost(), test.wantDefaultHost)
		}

		// Unset env variables after each test so they don't affect other tests
		err := os.Unsetenv(EnvDefaultHost)
		checkFatal(t, err)
		err = os.Unsetenv(EnvReposRoot)
		checkFatal(t, err)
	}
}
