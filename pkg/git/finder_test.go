package git

import (
	"git-get/pkg/git/test"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFinder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		reposMaker func(*testing.T) string
		want       int
	}{
		{
			name:       "single repos",
			reposMaker: makeSingleRepo,
			want:       1,
		}, {
			name:       "single nested repo",
			reposMaker: makeNestedRepo,
			want:       1,
		}, {
			name:       "multiple nested repo",
			reposMaker: makeMultipleNestedRepos,
			want:       2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			root := test.reposMaker(t)

			finder := NewRepoFinder(root)

			err := finder.Find()
			if err != nil {
				t.Fatalf("finder.Find() failed: %v", err)
			}

			assert.Len(t, finder.repos, test.want)
		})
	}
}

func TestExists(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		want error
	}{
		{
			name: "dir does not exist",
			path: "/this/directory/does/not/exist",
			want: ErrDirNotExist,
		}, {
			name: "dir exists",
			path: os.TempDir(),
			want: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := Exists(test.path)

			assert.ErrorIs(t, err, test.want)
		})
	}
}

func makeSingleRepo(t *testing.T) string {
	t.Helper()
	root := test.TempDir(t, "")

	test.RepoEmptyInDir(t, root)

	return root
}

func makeNestedRepo(t *testing.T) string {
	t.Helper()
	// a repo with single nested repo should still be counted as one beacause finder doesn't traverse inside nested repos
	root := test.TempDir(t, "")

	r := test.RepoEmptyInDir(t, root)
	test.RepoEmptyInDir(t, r.Path())

	return root
}

func makeMultipleNestedRepos(t *testing.T) string {
	t.Helper()
	root := test.TempDir(t, "")

	// create two repos inside root - should be counted as 2
	repo1 := test.RepoEmptyInDir(t, root)
	repo2 := test.RepoEmptyInDir(t, root)

	// created repos nested inside two parent roots - shouldn't be counted
	test.RepoEmptyInDir(t, repo1.Path())
	test.RepoEmptyInDir(t, repo1.Path())
	test.RepoEmptyInDir(t, repo2.Path())

	// create a empty dir inside root - shouldn't be counted
	test.TempDir(t, root)

	return root
}
