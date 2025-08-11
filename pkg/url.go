package pkg

import (
	"errors"
	"fmt"
	urlpkg "net/url"
	"path"
	"regexp"
	"strings"
)

var errEmptyURLPath = errors.New("parsed URL path is empty")

// scpSyntax matches the SCP-like addresses used by the ssh protocol (eg, [user@]host.xz:path/to/repo.git/).
// See: https://golang.org/src/cmd/go/internal/get/vcs.go
var scpSyntax = regexp.MustCompile(`^([a-zA-Z0-9_]+)@([a-zA-Z0-9._-]+):(.*)$`)

// ParseURL parses given rawURL string into a URL.
// When the parsed URL has an empty host, use the defaultHost.
// When the parsed URL has an empty scheme, use the defaultScheme.
func ParseURL(rawURL string, defaultHost string, defaultScheme string) (url *urlpkg.URL, err error) {
	// If rawURL matches the SCP-like syntax, convert it into a standard ssh Path.
	// eg, git@github.com:user/repo => ssh://git@github.com/user/repo
	if m := scpSyntax.FindStringSubmatch(rawURL); m != nil {
		url = &urlpkg.URL{
			Scheme: "ssh",
			User:   urlpkg.User(m[1]),
			Host:   m[2],
			Path:   path.Join("/", m[3]),
		}
	} else {
		url, err = urlpkg.Parse(rawURL)
		if err != nil {
			return nil, fmt.Errorf("failed parsing URL %s: %w", rawURL, err)
		}
	}

	if url.Host == "" && url.Path == "" {
		return nil, errEmptyURLPath
	}

	if url.Scheme == "git+ssh" {
		url.Scheme = "ssh"
	}

	// Default to configured defaultHost when host is empty
	if url.Host == "" {
		url.Host = defaultHost
		// Add a leading slash to path when host is missing. It's needed to correctly compare urlpkg.URL structs.
		url.Path = path.Join("/", url.Path)
	}

	// Default to configured defaultScheme when scheme is empty
	if url.Scheme == "" {
		url.Scheme = defaultScheme
	}

	// Default to "git" user when using ssh and no user is provided
	if url.Scheme == "ssh" && url.User == nil {
		url.User = urlpkg.User("git")
	}

	// Don't use host when scheme is file://. The fragment detected as url.Host should be a first directory of url.Path
	if url.Scheme == "file" && url.Host != "" {
		url.Path = path.Join(url.Host, url.Path)
		url.Host = ""
	}

	return url, nil
}

// URLToPath cleans up the URL and converts it into a path string.
// Eg, ssh://git@github.com:22/~user/repo.git => github.com/user/repo
//
// If skipHost is true, it removes the host part from the path.
// Eg, ssh://git@github.com:22/~user/repo.git => user/repo
func URLToPath(url urlpkg.URL, skipHost bool) string {
	// Remove port numbers from host.
	url.Host = strings.Split(url.Host, ":")[0]

	// Remove tilde (~) char from username.
	url.Path = strings.ReplaceAll(url.Path, "~", "")

	// Remove leading and trailing slashes.
	url.Path = strings.Trim(url.Path, "/")

	// Remove trailing ".git" from repo name.
	url.Path = strings.TrimSuffix(url.Path, ".git")

	if skipHost {
		return url.Path
	}

	if url.Host == "" {
		return url.Path
	}

	if url.Path == "" {
		return url.Host
	}

	return url.Host + "/" + url.Path
}
