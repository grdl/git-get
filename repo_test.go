package main

import "testing"

func TestFetch(t *testing.T) {
	// Create origin repo with a single commit in master
	origin := newTestRepo(t)
	createFile(t, origin, "file")
	stageFile(t, origin, "file")
	createCommit(t, origin, "Initial commit")

	// Clone the origin repo
	repo, err := CloneRepo(origin.Path(), newTempDir(t))
	checkFatal(t, err)

	// Check cloned status. It should not be behind origin
	status, err := NewRepoStatus(repo.Workdir())
	checkFatal(t, err)

	if status.BranchStatuses["master"].Behind != 0 {
		t.Errorf("Master should not be behind")
	}

	// Add another commit to origin
	createFile(t, origin, "anotherFile")
	stageFile(t, origin, "anotherFile")
	createCommit(t, origin, "Second commit")

	// Fetch cloned repo and check the status again
	err = Fetch(repo)
	status, err = NewRepoStatus(repo.Workdir())
	checkFatal(t, err)

	// Cloned master should now be 1 commit behind origin
	if status.BranchStatuses["master"].Behind != 1 {
		t.Errorf("Master should be 1 commit behind")
	}

	if status.BranchStatuses["master"].Ahead != 0 {
		t.Errorf("Master should not be ahead")
	}
}
