package main

import (
	"testing"

	"github.com/pkg/errors"
)

func TestNewBranch(t *testing.T) {
	repo := newTestRepo(t)

	createFile(t, repo, "file")
	stageFile(t, repo, "file")
	createCommit(t, repo, "Initial commit")
	branch := createBranch(t, repo, "branch")

	status, err := NewBranchStatus(branch)
	checkFatal(t, errors.Wrap(err, "Failed getting branch status"))

	if status.Name != "branch" {
		t.Errorf("Wrong branch name, got %s; want %s", status.Name, "branch")
	}

	if status.IsRemote != false {
		t.Errorf("Branch should be local")
	}
}
