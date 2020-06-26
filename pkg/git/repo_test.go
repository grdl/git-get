package git

import (
	"git-get/pkg/io"
	"git-get/pkg/test"
	"reflect"
	"testing"
)

func TestOpen(t *testing.T) {
	_, err := Open("/paththatdoesnotexist/repo")

	if err != io.ErrDirectoryAccess {
		t.Errorf("Opening a repo in non existing path should throw an error")
	}
}

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
			name:      "single tracked ",
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
			want:      0,
		},
		{
			name:      "single tracked ",
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

func TestCurrentBranch(t *testing.T) {
	tests := []struct {
		name      string
		repoMaker func(*testing.T) *test.Repo
		want      string
	}{
		// TODO: maybe add wantErr to check if error is returned correctly?
		// {
		// 	name:      "empty",
		// 	repoMaker: newTestRepo,
		// 	want:      "",
		// },
		{
			name:      "only master branch",
			repoMaker: test.RepoWithCommit,
			want:      master,
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
			name:      "only master branch",
			repoMaker: test.RepoWithCommit,
			want:      []string{"master"},
		},
		{
			name:      "new branch",
			repoMaker: test.RepoWithBranch,
			want:      []string{"feature/branch", "master"},
		},
		{
			name:      "checked out new tag",
			repoMaker: test.RepoWithTag,
			want:      []string{"master"},
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
			branch:    "master",
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
			name:      "master with upstream",
			repoMaker: test.RepoWithBranchWithUpstream,
			branch:    "master",
			want:      "origin/master",
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
			branch:    "master",
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
