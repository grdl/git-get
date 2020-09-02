package pkg

import (
	"git-get/pkg/cfg"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Following URLs are considered valid according to https://git-scm.com/docs/git-clone#_git_urls:
// ssh://[user@]host.xz[:port]/path/to/repo.git
// ssh://[user@]host.xz[:port]/~[user]/path/to/repo.git/
// [user@]host.xz:path/to/repo.git/
// [user@]host.xz:/~[user]/path/to/repo.git/
// git://host.xz[:port]/path/to/repo.git/
// git://host.xz[:port]/~[user]/path/to/repo.git/
// http[s]://host.xz[:port]/path/to/repo.git/
// ftp[s]://host.xz[:port]/path/to/repo.git/
// /path/to/repo.git/
// file:///path/to/repo.git/

func TestURLParse(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"ssh://github.com/grdl/git-get.git", "github.com/grdl/git-get"},
		{"ssh://user@github.com/grdl/git-get.git", "github.com/grdl/git-get"},
		{"ssh://user@github.com:1234/grdl/git-get.git", "github.com/grdl/git-get"},
		{"ssh://user@github.com/~user/grdl/git-get.git", "github.com/user/grdl/git-get"},
		{"git+ssh://github.com/grdl/git-get.git", "github.com/grdl/git-get"},
		{"git@github.com:grdl/git-get.git", "github.com/grdl/git-get"},
		{"git@github.com:/~user/grdl/git-get.git", "github.com/user/grdl/git-get"},
		{"git://github.com/grdl/git-get.git", "github.com/grdl/git-get"},
		{"git://github.com/~user/grdl/git-get.git", "github.com/user/grdl/git-get"},
		{"https://github.com/grdl/git-get.git", "github.com/grdl/git-get"},
		{"http://github.com/grdl/git-get.git", "github.com/grdl/git-get"},
		{"https://github.com/grdl/git-get", "github.com/grdl/git-get"},
		{"https://github.com/git-get.git", "github.com/git-get"},
		{"https://github.com/git-get", "github.com/git-get"},
		{"https://github.com/grdl/sub/path/git-get.git", "github.com/grdl/sub/path/git-get"},
		{"https://github.com:1234/grdl/git-get.git", "github.com/grdl/git-get"},
		{"https://github.com/grdl/git-get.git/", "github.com/grdl/git-get"},
		{"https://github.com/grdl/git-get/", "github.com/grdl/git-get"},
		{"https://github.com/grdl/git-get/////", "github.com/grdl/git-get"},
		{"https://github.com/grdl/git-get.git/////", "github.com/grdl/git-get"},
		{"ftp://github.com/grdl/git-get.git", "github.com/grdl/git-get"},
		{"ftps://github.com/grdl/git-get.git", "github.com/grdl/git-get"},
		{"rsync://github.com/grdl/git-get.git", "github.com/grdl/git-get"},
		{"local/grdl/git-get/", "github.com/local/grdl/git-get"},
		{"file://local/grdl/git-get", "local/grdl/git-get"},
	}

	for _, test := range tests {
		url, err := ParseURL(test.in, cfg.Defaults[cfg.KeyDefaultHost], cfg.Defaults[cfg.KeyDefaultScheme])
		assert.NoError(t, err)

		got := URLToPath(*url, false)
		assert.Equal(t, test.want, got)
	}
}
func TestURLParseSkipHost(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"ssh://github.com/grdl/git-get.git", "grdl/git-get"},
		{"ssh://user@github.com/grdl/git-get.git", "grdl/git-get"},
		{"ssh://user@github.com:1234/grdl/git-get.git", "grdl/git-get"},
		{"ssh://user@github.com/~user/grdl/git-get.git", "user/grdl/git-get"},
		{"git+ssh://github.com/grdl/git-get.git", "grdl/git-get"},
		{"git@github.com:grdl/git-get.git", "grdl/git-get"},
		{"git@github.com:/~user/grdl/git-get.git", "user/grdl/git-get"},
		{"git://github.com/grdl/git-get.git", "grdl/git-get"},
		{"git://github.com/~user/grdl/git-get.git", "user/grdl/git-get"},
		{"https://github.com/grdl/git-get.git", "grdl/git-get"},
		{"http://github.com/grdl/git-get.git", "grdl/git-get"},
		{"https://github.com/grdl/git-get", "grdl/git-get"},
		{"https://github.com/git-get.git", "git-get"},
		{"https://github.com/git-get", "git-get"},
		{"https://github.com/grdl/sub/path/git-get.git", "grdl/sub/path/git-get"},
		{"https://github.com:1234/grdl/git-get.git", "grdl/git-get"},
		{"https://github.com/grdl/git-get.git/", "grdl/git-get"},
		{"https://github.com/grdl/git-get/", "grdl/git-get"},
		{"https://github.com/grdl/git-get/////", "grdl/git-get"},
		{"https://github.com/grdl/git-get.git/////", "grdl/git-get"},
		{"ftp://github.com/grdl/git-get.git", "grdl/git-get"},
		{"ftps://github.com/grdl/git-get.git", "grdl/git-get"},
		{"rsync://github.com/grdl/git-get.git", "grdl/git-get"},
		{"local/grdl/git-get/", "local/grdl/git-get"},
		{"file://local/grdl/git-get", "local/grdl/git-get"},
	}

	for _, test := range tests {
		url, err := ParseURL(test.in, cfg.Defaults[cfg.KeyDefaultHost], cfg.Defaults[cfg.KeyDefaultScheme])
		assert.NoError(t, err)

		got := URLToPath(*url, true)
		assert.Equal(t, test.want, got)
	}
}

func TestDefaultScheme(t *testing.T) {
	tests := []struct {
		in     string
		scheme string
		want   string
	}{
		{"grdl/git-get", "ssh", "ssh://git@github.com/grdl/git-get"},
		{"grdl/git-get", "https", "https://github.com/grdl/git-get"},
		{"https://github.com/grdl/git-get", "ssh", "https://github.com/grdl/git-get"},
		{"https://github.com/grdl/git-get", "https", "https://github.com/grdl/git-get"},
		{"ssh://github.com/grdl/git-get", "ssh", "ssh://git@github.com/grdl/git-get"},
		{"ssh://github.com/grdl/git-get", "https", "ssh://git@github.com/grdl/git-get"},
		{"git+ssh://github.com/grdl/git-get", "https", "ssh://git@github.com/grdl/git-get"},
		{"git@github.com:grdl/git-get", "ssh", "ssh://git@github.com/grdl/git-get"},
		{"git@github.com:grdl/git-get", "https", "ssh://git@github.com/grdl/git-get"},
	}

	for _, test := range tests {
		url, err := ParseURL(test.in, cfg.Defaults[cfg.KeyDefaultHost], test.scheme)
		assert.NoError(t, err)

		want, err := url.Parse(test.want)
		assert.NoError(t, err)

		assert.Equal(t, url, want)
	}
}

func TestInvalidURLParse(t *testing.T) {
	invalidURLs := []string{
		"",
		//TODO: This Path is technically a correct scp-like syntax. Not sure how to handle it
		"github.com:grdl/git-git.get.git",

		//TODO: Is this a valid git Path?
		//"git@github.com:1234:grdl/git-get.git",
	}

	for _, test := range invalidURLs {
		_, err := ParseURL(test, cfg.Defaults[cfg.KeyDefaultHost], cfg.Defaults[cfg.KeyDefaultScheme])

		assert.Error(t, err)
	}
}
