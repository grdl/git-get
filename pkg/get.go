package pkg

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/grdl/git-get/pkg/git"
)

var ErrMissingRepoArg = errors.New("missing <REPO> argument or --dump flag")

// GetCfg provides configuration for the Get command.
type GetCfg struct {
	Branch    string
	DefHost   string
	DefScheme string
	Dump      string
	Root      string
	SkipHost  bool
	URL       string
}

// Get executes the "git get" command.
func Get(conf *GetCfg) error {
	if conf.URL == "" && conf.Dump == "" {
		return ErrMissingRepoArg
	}

	if conf.URL != "" {
		return cloneSingleRepo(conf)
	}

	if conf.Dump != "" {
		return cloneDumpFile(conf)
	}

	return nil
}

func cloneSingleRepo(conf *GetCfg) error {
	url, err := ParseURL(conf.URL, conf.DefHost, conf.DefScheme)
	if err != nil {
		return err
	}

	opts := &git.CloneOpts{
		URL:    url,
		Path:   filepath.Join(conf.Root, URLToPath(*url, conf.SkipHost)),
		Branch: conf.Branch,
	}

	_, err = git.Clone(opts)

	return err
}

func cloneDumpFile(conf *GetCfg) error {
	parsedLines, err := parseDumpFile(conf.Dump)
	if err != nil {
		return err
	}

	for _, line := range parsedLines {
		url, err := ParseURL(line.rawurl, conf.DefHost, conf.DefScheme)
		if err != nil {
			return err
		}

		opts := &git.CloneOpts{
			URL:    url,
			Path:   filepath.Join(conf.Root, URLToPath(*url, conf.SkipHost)),
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
