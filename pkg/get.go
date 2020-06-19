package pkg

import (
	"git-get/pkg/repo"
	"path"
)

// GetCfg provides configuration for the Get command.
type GetCfg struct {
	Branch  string
	DefHost string
	Dump    string
	Root    string
	URL     string
}

// Get executes the "git get" command.
func Get(c *GetCfg) error {
	// TODO: show something when no args

	if c.URL != "" {
		return cloneSingleRepo(c)
	}

	if c.Dump != "" {
		return cloneDumpFile(c)
	}
	return nil
}

func cloneSingleRepo(c *GetCfg) error {
	url, err := ParseURL(c.URL, c.DefHost)
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
	parsedLines, err := parseDumpFile(c.Dump)
	if err != nil {
		return err
	}

	for _, line := range parsedLines {
		url, err := ParseURL(line.rawurl, c.DefHost)
		if err != nil {
			return err
		}

		cloneOpts := &repo.CloneOpts{
			URL:            url,
			Path:           path.Join(c.Root, URLToPath(url)),
			Branch:         line.branch,
			IgnoreExisting: true,
		}

		_, err = repo.Clone(cloneOpts)
		if err != nil {
			return err
		}
	}
	return nil
}
