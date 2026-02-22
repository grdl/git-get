package pkg

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grdl/git-get/pkg/git"
)

var ErrMissingRepoArg = errors.New("missing <REPO> argument or --dump flag")

// GetCfg provides configuration for the Get command.
type GetCfg struct {
	Branch    string
	DefHost   string
	DefScheme string
	Dump      string
	Roots     []string
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

	pathSuffix := URLToPath(*url, conf.SkipHost)

	// Check if repo already exists in any of the roots
	for _, root := range conf.Roots {
		path := filepath.Join(root, pathSuffix)
		if exists, _ := git.Exists(path); exists {
			return nil
		}
	}

	// If not found, select most appropriate root
	root := selectBestRoot(conf.Roots, pathSuffix)
	opts := &git.CloneOpts{
		URL:    url,
		Path:   filepath.Join(root, pathSuffix),
		Branch: conf.Branch,
	}

	_, err = git.Clone(opts)

	return err
}

func selectBestRoot(roots []string, pathSuffix string) string {
	if len(roots) == 0 {
		return ""
	}

	// Heuristic: if host/owner directory present for similar repository
	// pathSuffix is e.g. "github.com/grdl/git-get"
	// We check if "github.com/grdl" exists in any root.
	// Or even just "github.com".

	parts := strings.Split(pathSuffix, "/")

	for i := len(parts) - 1; i > 0; i-- {
		prefix := filepath.Join(parts[:i]...)
		for _, root := range roots {
			if exists, _ := git.Exists(filepath.Join(root, prefix)); exists {
				return root
			}
		}
	}

	return roots[0]
}

// getRelativePaths returns a slice of paths relative to their respective roots.
func getRelativePaths(roots []string, paths []string) []string {
	relPaths := make([]string, len(paths))

	for idx, path := range paths {
		relPath := path

		for _, root := range roots {
			if strings.HasPrefix(path, root) {
				if p, err := filepath.Rel(root, path); err == nil {
					relPath = p

					break
				}
			}
		}

		relPaths[idx] = relPath
	}

	return relPaths
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

		pathSuffix := URLToPath(*url, conf.SkipHost)

		var targetPath string

		found := false

		// Check if repo already exists in any of the roots
		for _, root := range conf.Roots {
			path := filepath.Join(root, pathSuffix)
			if exists, _ := git.Exists(path); exists {
				targetPath = path
				found = true

				break
			}
		}

		if !found {
			root := selectBestRoot(conf.Roots, pathSuffix)
			targetPath = filepath.Join(root, pathSuffix)
		}

		opts := &git.CloneOpts{
			URL:    url,
			Path:   targetPath,
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
