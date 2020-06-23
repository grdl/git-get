package repo

import (
	"git-get/pkg/file"
	"reflect"
	"testing"
)

func TestOpen(t *testing.T) {
	_, err := Open("/paththatdoesnotexist/repo")

	if err != file.ErrDirectoryAccess {
		t.Errorf("Opening a repo in non existing path should throw an error")
	}
}

func TestUncommitted(t *testing.T) {
	tests := []struct {
		name      string
		repoMaker func(*testing.T) *testRepo
		want      int
	}{
		{
			name:      "empty",
			repoMaker: newTestRepo,
			want:      0,
		},
		{
			name:      "single untracked",
			repoMaker: testRepoWithUntracked,
			want:      0,
		},
		{
			name:      "single tracked ",
			repoMaker: testRepoWithStaged,
			want:      1,
		},
		{
			name:      "committed",
			repoMaker: testRepoWithCommit,
			want:      0,
		},
		{
			name:      "untracked and uncommitted",
			repoMaker: testRepoWithUncommittedAndUntracked,
			want:      1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := test.repoMaker(t)
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
		repoMaker func(*testing.T) *testRepo
		want      int
	}{
		{
			name:      "empty",
			repoMaker: newTestRepo,
			want:      0,
		},
		{
			name:      "single untracked",
			repoMaker: testRepoWithUntracked,
			want:      0,
		},
		{
			name:      "single tracked ",
			repoMaker: testRepoWithStaged,
			want:      1,
		},
		{
			name:      "committed",
			repoMaker: testRepoWithCommit,
			want:      0,
		},
		{
			name:      "untracked and uncommitted",
			repoMaker: testRepoWithUncommittedAndUntracked,
			want:      1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := test.repoMaker(t)
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
		repoMaker func(*testing.T) *testRepo
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
			repoMaker: testRepoWithCommit,
			want:      master,
		},
		{
			name:      "checked out new branch",
			repoMaker: testRepoWithBranch,
			want:      "feature/branch",
		},
		{
			name:      "checked out new tag",
			repoMaker: testRepoWithTag,
			want:      head,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := test.repoMaker(t)
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
		repoMaker func(*testing.T) *testRepo
		want      []string
	}{
		{
			name:      "empty",
			repoMaker: newTestRepo,
			want:      []string{""},
		},
		{
			name:      "only master branch",
			repoMaker: testRepoWithCommit,
			want:      []string{"master"},
		},
		{
			name:      "new branch",
			repoMaker: testRepoWithBranch,
			want:      []string{"feature/branch", "master"},
		},
		{
			name:      "checked out new tag",
			repoMaker: testRepoWithTag,
			want:      []string{"master"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := test.repoMaker(t)
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
		repoMaker func(*testing.T) *testRepo
		branch    string
		want      string
	}{
		{
			name:      "empty",
			repoMaker: newTestRepo,
			branch:    "master",
			want:      "",
		},
		// TODO: add wantErr
		{
			name:      "wrong branch name",
			repoMaker: testRepoWithCommit,
			branch:    "wrong_branch_name",
			want:      "",
		},
		{
			name:      "master with upstream",
			repoMaker: testRepoWithBranchWithUpstream,
			branch:    "master",
			want:      "origin/master",
		},
		{
			name:      "branch with upstream",
			repoMaker: testRepoWithBranchWithUpstream,
			branch:    "feature/branch",
			want:      "origin/feature/branch",
		},
		{
			name:      "branch without upstream",
			repoMaker: testRepoWithBranchWithoutUpstream,
			branch:    "feature/branch",
			want:      "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := test.repoMaker(t)
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
		repoMaker func(*testing.T) *testRepo
		branch    string
		want      []int
	}{
		{
			name:      "fresh clone",
			repoMaker: testRepoWithBranchWithUpstream,
			branch:    "master",
			want:      []int{0, 0},
		},
		{
			name:      "branch ahead",
			repoMaker: testRepoWithBranchAhead,
			branch:    "feature/branch",
			want:      []int{1, 0},
		},

		{
			name:      "branch behind",
			repoMaker: testRepoWithBranchBehind,
			branch:    "feature/branch",
			want:      []int{0, 1},
		},
		{
			name:      "branch ahead and behind",
			repoMaker: testRepoWithBranchAheadAndBehind,
			branch:    "feature/branch",
			want:      []int{2, 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := test.repoMaker(t)
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

// func TestClone(t *testing.T) {
// 	url, _ := url.Parse("https://github.com/grdl/pronestheus")
// 	opts := &CloneOpts{
// 		URL:  url,
// 		Path: "/tmp/stuff/nanana",
// 	}

// 	repo, err := Clone(opts)
// 	if err != nil {
// 		t.Errorf("got error %q", err)
// 	}

// 	fmt.Println(repo.path)
// }
