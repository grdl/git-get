package pkg

import (
	"git-get/pkg/repo"
	"path"
)

// GetCfg provides configuration for the Get command.
type GetCfg struct {
	Branch string
	Dump   string
	Root   string
	URL    string
}

// Get executes the "git get" command.
func Get(c *GetCfg) error {
	if c.Dump != "" {
		return cloneDumpFile(c)
	}

	if c.URL != "" {
		return cloneSingleRepo(c)
	}

	return nil
}

func cloneSingleRepo(c *GetCfg) error {
	url, err := ParseURL(c.URL)
	if err != nil {
		return err
	}

	cloneOpts := &repo.CloneOpts{
		URL:    url,
		Path:   path.Join(c.Root, URLToPath(url)),
		Branch: c.Branch,
	}

	_, err = repo.Clone(cloneOpts)

	return err
}

func cloneDumpFile(c *GetCfg) error {
	opts, err := ParseDumpFile(c.Dump)
	if err != nil {
		return err
	}

	for _, opt := range opts {
		path := path.Join(c.Root, URLToPath(opt.URL))
		opt.Path = path
		_, _ = repo.Clone(opt)
	}
	return nil
}
