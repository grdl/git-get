package git

import (
	"git-get/pkg/io"
	"git-get/pkg/run"
	"git-get/pkg/test"
	"path/filepath"
	"testing"
)

// cfgStub represents a gitconfig file but instead of using a global one, it creates a temporary git repo and uses its local gitconfig.
type cfgStub struct {
	repo *test.Repo
}

func newCfgStub(t *testing.T) *cfgStub {
	r := test.RepoEmpty(t)
	return &cfgStub{
		repo: r,
	}
}

func (c *cfgStub) Get(key string) string {
	out, err := run.Git("config", "--local", key).OnRepo(c.repo.Path()).AndCaptureLine()
	if err != nil {
		return ""
	}

	return out
}

func TestGitConfig(t *testing.T) {
	tests := []struct {
		name        string
		configMaker func(t *testing.T) *cfgStub
		key         string
		want        string
	}{
		{
			name:        "empty",
			configMaker: makeConfigEmpty,
			key:         "gitget.host",
			want:        "",
		},
		{
			name:        "valid",
			configMaker: makeConfigValid,
			key:         "gitget.host",
			want:        "github.com",
		}, {
			name:        "only section name",
			configMaker: makeConfigValid,
			key:         "gitget",
			want:        "",
		}, {
			name:        "missing key",
			configMaker: makeConfigValid,
			key:         "gitget.missingkey",
			want:        "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := test.configMaker(t)

			got := cfg.Get(test.key)

			if got != test.want {
				t.Errorf("expected %q; got %q", test.want, got)
			}
		})
	}
}

func makeConfigEmpty(t *testing.T) *cfgStub {
	c := newCfgStub(t)
	io.Write(filepath.Join(c.repo.Path(), dotgit, "config"), "")

	return c
}

func makeConfigValid(t *testing.T) *cfgStub {
	c := newCfgStub(t)

	gitconfig := `
	[user]
		name = grdl
	[gitget]
		host = github.com
	`
	io.Write(filepath.Join(c.repo.Path(), dotgit, "config"), gitconfig)

	return c
}
