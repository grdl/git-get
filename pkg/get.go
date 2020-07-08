package pkg

import (
	"fmt"
	"git-get/pkg/git"
	"path/filepath"
)

// GetCfg provides configuration for the Get command.
type GetCfg struct {
	Branch   string
	DefHost  string
	Dump     string
	Root     string
	SkipHost bool
	URL      string
}

// Get executes the "git get" command.
func Get(c *GetCfg) error {
	if c.URL == "" && c.Dump == "" {
		return fmt.Errorf("missing <REPO> argument or --dump flag")
	}

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

	opts := &git.CloneOpts{
		URL:    url,
		Path:   filepath.Join(c.Root, URLToPath(*url, c.SkipHost)),
		Branch: c.Branch,
	}

	_, err = git.Clone(opts)

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

		opts := &git.CloneOpts{
			URL:    url,
			Path:   filepath.Join(c.Root, URLToPath(*url, c.SkipHost)),
			Branch: line.branch,
		}

		// If target path already exists, skip cloning this repo
		if exists, _ := git.Exists(opts.Path); exists {
			continue
		}

		fmt.Printf("Cloning %s...\n", opts.URL.String())
		_, err = git.Clone(opts)
		if err != nil {
			return err
		}
	}
	return nil
}
