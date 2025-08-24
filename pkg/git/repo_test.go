package git

import (
	"git-get/pkg/git/test"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUncommitted(t *testing.T) {
	tests := []struct {
		name      string
		repoMaker func(*testing.T) *test.Repo
		want      int
	}{
		{
			name:      "empty",
			repoMaker: test.RepoEmpty,
			want:      0,
		},
		{
			name:      "single untracked",
			repoMaker: test.RepoWithUntracked,
			want:      0,
		},
		{
			name:      "single tracked",
			repoMaker: test.RepoWithStaged,
			want:      1,
		},
		{
			name:      "committed",
			repoMaker: test.RepoWithCommit,
			want:      0,
		},
		{
			name:      "untracked and uncommitted",
			repoMaker: test.RepoWithUncommittedAndUntracked,
			want:      1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, _ := Open(test.repoMaker(t).Path())

			got, err := r.Uncommitted()
			if err != nil {
				t.Errorf("got error %q", err)
			}

			if got != test.want {
				t.Errorf("expected %d; got %d", test.want, got)
			}
		})
	}
}
func TestUntracked(t *testing.T) {
	tests := []struct {
		name      string
		repoMaker func(*testing.T) *test.Repo
		want      int
	}{
		{
			name:      "empty",
			repoMaker: test.RepoEmpty,
			want:      0,
		},
		{
			name:      "single untracked",
			repoMaker: test.RepoWithUntracked,
			want:      1,
		},
		{
			name:      "single tracked ",
			repoMaker: test.RepoWithStaged,
			want:      0,
		},
		{
			name:      "committed",
			repoMaker: test.RepoWithCommit,
			want:      0,
		},
		{
			name:      "untracked and uncommitted",
			repoMaker: test.RepoWithUncommittedAndUntracked,
			want:      1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, _ := Open(test.repoMaker(t).Path())

			got, err := r.Untracked()
			if err != nil {
				t.Errorf("got error %q", err)
			}

			if got != test.want {
				t.Errorf("expected %d; got %d", test.want, got)
			}
		})
	}
}

func TestCurrentBranch(t *testing.T) {
	tests := []struct {
		name      string
		repoMaker func(*testing.T) *test.Repo
		want      string
	}{
		{
			name:      "empty repo without commits",
			repoMaker: test.RepoEmpty,
			want:      "main",
		},
		{
			name:      "only main branch",
			repoMaker: test.RepoWithCommit,
			want:      main,
		},
		{
			name:      "checked out new branch",
			repoMaker: test.RepoWithBranch,
			want:      "feature/branch",
		},
		{
			name:      "checked out new tag",
			repoMaker: test.RepoWithTag,
			want:      head,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, _ := Open(test.repoMaker(t).Path())

			got, err := r.CurrentBranch()
			if err != nil {
				t.Errorf("got error %q", err)
			}

			if got != test.want {
				t.Errorf("expected %q; got %q", test.want, got)
			}
		})
	}
}
func TestBranches(t *testing.T) {
	tests := []struct {
		name      string
		repoMaker func(*testing.T) *test.Repo
		want      []string
	}{
		{
			name:      "empty",
			repoMaker: test.RepoEmpty,
			want:      []string{""},
		},
		{
			name:      "only main branch",
			repoMaker: test.RepoWithCommit,
			want:      []string{"main"},
		},
		{
			name:      "new branch",
			repoMaker: test.RepoWithBranch,
			want:      []string{"feature/branch", "main"},
		},
		{
			name:      "checked out new tag",
			repoMaker: test.RepoWithTag,
			want:      []string{"main"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, _ := Open(test.repoMaker(t).Path())

			got, err := r.Branches()
			if err != nil {
				t.Errorf("got error %q", err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("expected %+v; got %+v", test.want, got)
			}
		})
	}
}
func TestUpstream(t *testing.T) {
	tests := []struct {
		name      string
		repoMaker func(*testing.T) *test.Repo
		branch    string
		want      string
	}{
		{
			name:      "empty",
			repoMaker: test.RepoEmpty,
			branch:    "main",
			want:      "",
		},
		// TODO: add wantErr
		{
			name:      "wrong branch name",
			repoMaker: test.RepoWithCommit,
			branch:    "wrong_branch_name",
			want:      "",
		},
		{
			name:      "main with upstream",
			repoMaker: test.RepoWithBranchWithUpstream,
			branch:    "main",
			want:      "origin/main",
		},
		{
			name:      "branch with upstream",
			repoMaker: test.RepoWithBranchWithUpstream,
			branch:    "feature/branch",
			want:      "origin/feature/branch",
		},
		{
			name:      "branch without upstream",
			repoMaker: test.RepoWithBranchWithoutUpstream,
			branch:    "feature/branch",
			want:      "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, _ := Open(test.repoMaker(t).Path())
			got, _ := r.Upstream(test.branch)

			// TODO:
			// if err != nil {
			// 	t.Errorf("got error %q", err)
			// }

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("expected %+v; got %+v", test.want, got)
			}
		})
	}
}
func TestAheadBehind(t *testing.T) {
	tests := []struct {
		name      string
		repoMaker func(*testing.T) *test.Repo
		branch    string
		want      []int
	}{
		{
			name:      "fresh clone",
			repoMaker: test.RepoWithBranchWithUpstream,
			branch:    "main",
			want:      []int{0, 0},
		},
		{
			name:      "branch ahead",
			repoMaker: test.RepoWithBranchAhead,
			branch:    "feature/branch",
			want:      []int{1, 0},
		},

		{
			name:      "branch behind",
			repoMaker: test.RepoWithBranchBehind,
			branch:    "feature/branch",
			want:      []int{0, 1},
		},
		{
			name:      "branch ahead and behind",
			repoMaker: test.RepoWithBranchAheadAndBehind,
			branch:    "feature/branch",
			want:      []int{2, 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, _ := Open(test.repoMaker(t).Path())

			upstream, err := r.Upstream(test.branch)
			if err != nil {
				t.Errorf("got error %q", err)
			}

			ahead, behind, err := r.AheadBehind(test.branch, upstream)
			if err != nil {
				t.Errorf("got error %q", err)
			}

			if ahead != test.want[0] || behind != test.want[1] {
				t.Errorf("expected %+v; got [%d, %d]", test.want, ahead, behind)
			}
		})
	}
}

func TestCleanupFailedClone(t *testing.T) {
	// Test dir structure:
	// root
	// └── a/
	//     ├── b/
	//     │   └── c/
	//     └── x/
	//         └── y/
	//        	   └── file.txt
	tests := []struct {
		path     string // path to cleanup
		wantGone string // this path should be deleted, if empty - nothing should be deleted
		wantStay string // this path shouldn't be deleted
	}{
		{
			path:     "a/b/c/repo",
			wantGone: "a/b/c/repo",
			wantStay: "a",
		}, {
			path:     "a/b/c/repo",
			wantGone: "a/b",
			wantStay: "a",
		}, {
			path:     "a/b/repo",
			wantGone: "",
			wantStay: "a/b/c",
		}, {
			path:     "a/x/y/repo",
			wantGone: "",
			wantStay: "a/x/y",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			root := createTestDirTree(t)

			path := filepath.Join(root, test.path)
			cleanupFailedClone(path)

			if test.wantGone != "" {
				wantGone := filepath.Join(root, test.wantGone)
				assert.NoDirExists(t, wantGone, "%s dir should be deleted during the cleanup", wantGone)
			}

			if test.wantStay != "" {
				wantLeft := filepath.Join(root, test.wantStay)
				assert.DirExists(t, wantLeft, "%s dir should not be deleted during the cleanup", wantLeft)
			}
		})
	}
}

func TestRemote(t *testing.T) {
	tests := []struct {
		name      string
		repoMaker func(*testing.T) *test.Repo
		want      string
		wantErr   bool
	}{
		{
			name:      "empty repo without remote",
			repoMaker: test.RepoEmpty,
			want:      "",
			wantErr:   false,
		},
		{
			name:      "repo with commit but no remote",
			repoMaker: test.RepoWithCommit,
			want:      "",
			wantErr:   false,
		},
		{
			name:      "repo with upstream",
			repoMaker: test.RepoWithBranchWithUpstream,
			want:      "", // This will contain the actual remote URL but we just test it doesn't error
			wantErr:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, _ := Open(test.repoMaker(t).Path())
			got, err := r.Remote()

			if test.wantErr && err == nil {
				t.Errorf("expected error but got none")
			}

			if !test.wantErr && err != nil {
				t.Errorf("unexpected error: %q", err)
			}

			// For repos with remote, just check no error occurred
			if test.name == "repo with upstream" {
				if err != nil {
					t.Errorf("unexpected error for repo with remote: %q", err)
				}
			} else if got != test.want {
				t.Errorf("expected %q; got %q", test.want, got)
			}
		})
	}
}

func createTestDirTree(t *testing.T) string {
	t.Helper()
	root := test.TempDir(t, "")

	err := os.MkdirAll(filepath.Join(root, "a", "b", "c"), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(filepath.Join(root, "a", "x", "y"), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Create(filepath.Join(root, "a", "x", "y", "file.txt"))
	if err != nil {
		t.Fatal(err)
	}

	return root
}
