package git

import (
	"testing"
)

// cfgStub represents a gitconfig file but instead of using a global one, it creates a temporary git repo and uses its local gitconfig.
type cfgStub struct {
	repo *testRepo
}

func newCfgStub(t *testing.T) *cfgStub {
	r := testRepoEmpty(t)
	return &cfgStub{
		repo: r,
	}
}

func (c *cfgStub) Get(key string) string {
	cmd := gitCmd(c.repo.path, "config", "--local", key)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := lines(out)
	return lines[0]
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
	c.repo.writeFile(".git/config", "")

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
	c.repo.writeFile(".git/config", gitconfig)

	return c
}
