package pkg

import "testing"

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
		{"local/grdl/git-get/", "local/grdl/git-get"},
		{"file://local/grdl/git-get", "local/grdl/git-get"},
	}

	for _, test := range tests {
		got, err := URLToPath(test.in)
		if err != nil {
			t.Errorf("Error parsing URL: %+v", err)
		}

		if got != test.want {
			t.Errorf("Wrong result of parsing URL: %s, got: %s; want: %s", test.in, got, test.want)
		}
	}
}

func TestInvalidURLParse(t *testing.T) {
	invalidURLs := []string{
		"",
		//TODO: This URL is technically a correct scp-like syntax. Not sure how to handle it
		"github.com:grdl/git-git.get.git",
		"git@github.com:1234:grdl/git-get.git",
	}

	for _, url := range invalidURLs {
		got, err := URLToPath(url)
		if err == nil {
			t.Errorf("Wrong result of parsing invalid URL: %s, got: %s, want: error", url, got)
		}
	}

}
