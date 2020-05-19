package main

import (
	"reflect"
	"testing"
)

func TestNewLocalBranch(t *testing.T) {
	repo := newTestRepo(t)

	createFile(t, repo, "file")
	stageFile(t, repo, "file")
	createCommit(t, repo, "Initial commit")
	branch := createBranch(t, repo, "branch")

	status, err := NewBranchStatus(repo, branch)
	checkFatal(t, err)

	want := BranchStatus{
		Name:        "branch",
		IsRemote:    false,
		HasUpstream: false,
		NeedsPull:   false,
		NeedsPush:   false,
		Ahead:       0,
		Behind:      0,
	}

	if status != want {
		t.Errorf("Wrong branch status, got %+v; want %+v", status, want)
	}
}

func TestClonedBranches(t *testing.T) {
	origin := newTestRepo(t)
	createFile(t, origin, "file")
	stageFile(t, origin, "file")
	createCommit(t, origin, "Initial commit")

	repo, err := CloneRepo(origin.Path(), newTempDir(t))
	checkFatal(t, err)

	createBranch(t, repo, "branch")

	branches, err := Branches(repo)
	checkFatal(t, err)

	var tests = []struct {
		got  BranchStatus
		want BranchStatus
	}{
		{branches["master"], BranchStatus{
			Name:        "master",
			IsRemote:    false,
			HasUpstream: true,
		}},
		{branches["origin/master"], BranchStatus{
			Name:        "origin/master",
			IsRemote:    true,
			HasUpstream: false,
		}},
		{branches["branch"], BranchStatus{
			Name:        "branch",
			IsRemote:    false,
			HasUpstream: false,
		}},
	}

	for _, test := range tests {
		if !reflect.DeepEqual(test.got, test.want) {
			t.Errorf("Wrong branch status, got %+v; want %+v", test.got, test.want)
		}
	}
}
