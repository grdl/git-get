package pkg

import (
	"reflect"
	"testing"
)

func TestStatus(t *testing.T) {
	var tests = []struct {
		makeTestRepo func(*testing.T) *Repo
		want         *RepoStatus
	}{
		{newRepoEmpty, &RepoStatus{
			HasUntrackedFiles:     false,
			HasUncommittedChanges: false,
			Branches:              nil,
		}},
		{newRepoWithUntracked, &RepoStatus{
			HasUntrackedFiles:     true,
			HasUncommittedChanges: false,
			Branches:              nil,
		}},
		{newRepoWithStaged, &RepoStatus{
			HasUntrackedFiles:     false,
			HasUncommittedChanges: true,
			Branches:              nil,
		}},
		{newRepoWithCommit, &RepoStatus{
			HasUntrackedFiles:     false,
			HasUncommittedChanges: false,
			Branches: []*BranchStatus{
				{
					Name:      "master",
					Upstream:  "",
					NeedsPull: false,
					NeedsPush: false,
				},
			},
		}},
		{newRepoWithModified, &RepoStatus{
			HasUntrackedFiles:     false,
			HasUncommittedChanges: true,
			Branches: []*BranchStatus{
				{
					Name:      "master",
					Upstream:  "",
					NeedsPull: false,
					NeedsPush: false,
				},
			},
		}},
		{newRepoWithIgnored, &RepoStatus{
			HasUntrackedFiles:     false,
			HasUncommittedChanges: false,
			Branches: []*BranchStatus{
				{
					Name:      "master",
					Upstream:  "",
					NeedsPull: false,
					NeedsPush: false,
				},
			},
		}},
		{newRepoWithLocalBranch, &RepoStatus{
			HasUntrackedFiles:     false,
			HasUncommittedChanges: false,
			Branches: []*BranchStatus{
				{
					Name:      "local",
					Upstream:  "",
					NeedsPull: false,
					NeedsPush: false,
				}, {
					Name:      "master",
					Upstream:  "",
					NeedsPull: false,
					NeedsPush: false,
				},
			},
		}},
		{newRepoWithClonedBranch, &RepoStatus{
			HasUntrackedFiles:     false,
			HasUncommittedChanges: false,
			Branches: []*BranchStatus{
				{
					Name:      "local",
					Upstream:  "",
					NeedsPull: false,
					NeedsPush: false,
				}, {
					Name:      "master",
					Upstream:  "origin/master",
					NeedsPull: false,
					NeedsPush: false,
				},
			},
		}},
		{newRepoWithBranchAhead, &RepoStatus{
			HasUntrackedFiles:     false,
			HasUncommittedChanges: false,
			Branches: []*BranchStatus{
				{
					Name:      "master",
					Upstream:  "origin/master",
					NeedsPull: false,
					NeedsPush: true,
				},
			},
		}},
		{newRepoWithBranchBehind, &RepoStatus{
			HasUntrackedFiles:     false,
			HasUncommittedChanges: false,
			Branches: []*BranchStatus{
				{
					Name:      "master",
					Upstream:  "origin/master",
					NeedsPull: true,
					NeedsPush: false,
				},
			},
		}},
		{newRepoWithBranchAheadAndBehind, &RepoStatus{
			HasUntrackedFiles:     false,
			HasUncommittedChanges: false,
			Branches: []*BranchStatus{
				{
					Name:      "master",
					Upstream:  "origin/master",
					NeedsPull: true,
					NeedsPush: true,
				},
			},
		}},
	}

	for _, test := range tests {
		repo := test.makeTestRepo(t)

		err := repo.LoadStatus()
		checkFatal(t, err)

		if !reflect.DeepEqual(repo.Status, test.want) {
			t.Errorf("Wrong repo status, got: %+v; want: %+v", repo.Status, test.want)
		}
	}
}

// TODO: test branch status when tracking a local branch
