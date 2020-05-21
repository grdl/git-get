package pkg

import (
	"net/url"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// TODO: maybe use https://github.com/whilp/git-urls?

func URLToPath(rawurl string) (string, error) {
	parsed, err := url.Parse(rawurl)
	if err != nil {
		return "", errors.Wrap(err, "Failed parsing URL")
	}

	repoHost := strings.Split(parsed.Host, ":")[0]

	repoPath := parsed.Path
	//repoPath = strings.TrimSuffix(repoPath, ".git/")
	//repoPath = strings.TrimSuffix(repoPath, ".git")

	localPath := path.Join(repoHost, repoPath)
	localPath = strings.TrimSuffix(localPath, ".git")
	return localPath, nil
}
