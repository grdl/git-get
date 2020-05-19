package main

import (
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
