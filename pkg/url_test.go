package pkg

import (
	"git-get/pkg/cfg"
	"testing"
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

	// We need to init config first so the default values are correctly loaded
	cfg.Init()

	for _, test := range tests {
		url, err := ParseURL(test.in)
		if err != nil {
			t.Errorf("Error parsing Path: %+v", err)
		}

		got := URLToPath(url)

		if got != test.want {
			t.Errorf("Wrong result of parsing Path: %s, got: %s; wantBranch: %s", test.in, got, test.want)
		}
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

	for _, in := range invalidURLs {
		got, err := ParseURL(in)
		if err == nil {
			t.Errorf("Wrong result of parsing invalid Path: %s, got: %s, wantBranch: error", in, got)
		}
	}
}
