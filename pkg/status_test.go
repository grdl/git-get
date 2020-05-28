package pkg

import (
	"reflect"
	"testing"
)

func TestStatus(t *testing.T) {
	var tests = []struct {
		makeTestRepo func(*testing.T) *TestRepo
		want         *RepoStatus
	}{
		{NewRepoEmpty, &RepoStatus{
			HasUntrackedFiles:     false,
			HasUncommittedChanges: false,
			Branches:              nil,
		}},
		{NewRepoWithUntracked, &RepoStatus{
			HasUntrackedFiles:     true,
			HasUncommittedChanges: false,
			Branches:              nil,
		}},
		{NewRepoWithStaged, &RepoStatus{
			HasUntrackedFiles:     false,
			HasUncommittedChanges: true,
			Branches:              nil,
		}},
		{NewRepoWithCommit, &RepoStatus{
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
		{NewRepoWithModified, &RepoStatus{
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
		{NewRepoWithIgnored, &RepoStatus{
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
		{NewRepoWithLocalBranch, &RepoStatus{
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
		{NewRepoWithClonedBranch, &RepoStatus{
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
		{NewRepoWithBranchAhead, &RepoStatus{
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
		{NewRepoWithBranchBehind, &RepoStatus{
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
		{NewRepoWithBranchAheadAndBehind, &RepoStatus{
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
		tr := test.makeTestRepo(t)

		repo, err := OpenRepo(tr.Path)
		checkFatal(t, err)

		err = repo.LoadStatus()
		checkFatal(t, err)

		if !reflect.DeepEqual(repo.Status, test.want) {
			t.Errorf("Wrong repo status, got: %+v; want: %+v", repo.Status, test.want)
		}
	}
}
