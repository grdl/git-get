package path

import (
	"git-get/cfg"
	urlpkg "net/url"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// scpSyntax matches the SCP-like addresses used by the ssh protocol (eg, [user@]host.xz:path/to/repo.git/).
// See: https://golang.org/src/cmd/go/internal/get/vcs.go
var scpSyntax = regexp.MustCompile(`^([a-zA-Z0-9_]+)@([a-zA-Z0-9._-]+):(.*)$`)

func ParseURL(rawURL string) (url *urlpkg.URL, err error) {
	// If rawURL matches the SCP-like syntax, convert it into a standard ssh Path.
	// eg, git@github.com:user/repo => ssh://git@github.com/user/repo
	if m := scpSyntax.FindStringSubmatch(rawURL); m != nil {
		url = &urlpkg.URL{
			Scheme: "ssh",
			User:   urlpkg.User(m[1]),
			Host:   m[2],
			Path:   m[3],
		}
	} else {
		url, err = urlpkg.Parse(rawURL)
		if err != nil {
			return nil, errors.Wrap(err, "Failed parsing Path")
		}
	}

	if url.Host == "" && url.Path == "" {
		return nil, errors.New("Parsed Path is empty")
	}

	if url.Scheme == "git+ssh" {
		url.Scheme = "ssh"
	}

	// Default to "git" user when using ssh and no user is provided
	if url.Scheme == "ssh" && url.User == nil {
		url.User = urlpkg.User("git")
	}

	// Default to configured defaultHost when host is empty
	if url.Host == "" {
		url.Host = viper.GetString(cfg.KeyDefaultHost)
	}

	// Default to https when scheme is empty
	if url.Scheme == "" {
		url.Scheme = "https"
	}

	return url, nil
}

func URLToPath(url *urlpkg.URL) (repoPath string) {
	// Remove port numbers from host
	repoHost := strings.Split(url.Host, ":")[0]

	// Remove trailing ".git" from repo name
	repoPath = path.Join(repoHost, url.Path)
	repoPath = strings.TrimSuffix(repoPath, ".git")

	// Remove tilde (~) char from username
	repoPath = strings.ReplaceAll(repoPath, "~", "")

	return repoPath
}
