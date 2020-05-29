package pkg

import (
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5/config"
	plumbing "github.com/go-git/go-git/v5/plumbing/format/config"
	"github.com/mitchellh/go-homedir"
)

const (
	CfgSection     = "gitget"
	CfgReposRoot   = "reposRoot"
	CfgDefaultHost = "defaultHost"

	EnvReposRoot   = "GITGET_REPOSROOT"
	EnvDefaultHost = "GITGET_DEFAULTHOST"

	DefaultReposRootSubpath = "repositories"
	DefaultHost             = "github.com"
)

var Cfg *Conf

// Conf provides methods for accessing configuration values
// Values are looked up in the following order: env variable, gitignore file, default value
type Conf struct {
	gitconfig *config.Config
}

func LoadConf() {
	// We don't care if loading gitconfig file throws an error
	// When gitconfig is nil, getters will just return the default values
	gitconfig, _ := config.LoadConfig(config.GlobalScope)

	Cfg = &Conf{
		gitconfig: gitconfig,
	}
}

func (c *Conf) ReposRoot() string {
	defReposRoot := path.Join(home(), DefaultReposRootSubpath)

	reposRoot := os.Getenv(EnvReposRoot)
	if reposRoot != "" {
		return reposRoot
	}

	if c.gitconfig == nil {
		return defReposRoot
	}

	gitget := c.findConfigSection(CfgSection)
	if gitget == nil {
		return defReposRoot
	}

	reposRoot = gitget.Option(CfgReposRoot)
	reposRoot = strings.TrimSpace(reposRoot)
	if reposRoot == "" {
		return defReposRoot
	}

	return reposRoot
}

func (c *Conf) DefaultHost() string {
	defaultHost := os.Getenv(EnvDefaultHost)
	if defaultHost != "" {
		return defaultHost
	}

	if c.gitconfig == nil {
		return DefaultHost
	}

	gitget := c.findConfigSection(CfgSection)
	if gitget == nil {
		return DefaultHost
	}

	defaultHost = gitget.Option(CfgDefaultHost)
	defaultHost = strings.TrimSpace(defaultHost)
	if defaultHost == "" {
		return DefaultHost
	}

	return defaultHost
}

func (c *Conf) findConfigSection(name string) *plumbing.Section {
	for _, s := range c.gitconfig.Raw.Sections {
		if s.Name == name {
			return s
		}
	}

	return nil
}

// home returns path to a home directory or empty string if can't be found
// Using empty string means that in the unlikely situation where home dir can't be found
// and there's no reposRoot specified in the global config, the current dir will be used as repos root.
func home() string {
	home, err := homedir.Dir()
	if err != nil {
		return ""
	}

	return home
}
