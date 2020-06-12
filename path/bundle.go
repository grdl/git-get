package path

import (
	"bufio"
	"git-get/git"
	"os"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrInvalidNumberOfElements = errors.New("More than two space-separated 2 elements on the line")
)

// ParseBundleFile opens a given gitgetfile and parses its content into a slice of CloneOpts.
func ParseBundleFile(path string) ([]*git.CloneOpts, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed opening gitgetfile %s", path)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var opts []*git.CloneOpts
	var line int
	for scanner.Scan() {
		line++
		opt, err := parseLine(scanner.Text())
		if err != nil {
			return nil, errors.Wrapf(err, "Failed parsing line %d", line)
		}

		opts = append(opts, opt)
	}

	return opts, nil
}

// parseLine splits a gitgetfile line into space-separated segments.
// First part is the URL to clone. Second, optional, is the branch (or tag) to checkout after cloning
func parseLine(line string) (*git.CloneOpts, error) {
	parts := strings.Split(line, " ")

	if len(parts) > 2 {
		return nil, ErrInvalidNumberOfElements
	}

	url, err := ParseURL(parts[0])
	if err != nil {
		return nil, err
	}

	branch := ""
	if len(parts) == 2 {
		branch = parts[1]
	}

	return &git.CloneOpts{
		URL:    url,
		Branch: branch,
		// When cloning a bundle we ignore errors about already cloned repos
		IgnoreExisting: true,
	}, nil
}
