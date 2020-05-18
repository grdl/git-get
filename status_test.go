package main

import (
	"testing"

	git "github.com/libgit2/git2go/v30"
)

func TestStatusWithEmptyRepo(t *testing.T) {
	repo, err := newTestRepo()
	checkFatal(t, err)
	defer cleanupRepo(t, repo)

	entries, err := statusEntries(repo)
	checkFatal(t, err)

	if len(entries) != 0 {
		t.Errorf("Empty repo should have no status entries")
	}

	statuses, err := GetStatus(repo.Workdir())
	checkFatal(t, err)

	if len(statuses) != 1 && statuses[0] != StatusOk {
		t.Errorf("Empty repo should have a single StatusOk status")
	}
}

func TestStatusWithSingleUnstagedFile(t *testing.T) {
	repo, err := newTestRepo()
	checkFatal(t, err)
	defer cleanupRepo(t, repo)

	err = createFile(repo, "SomeFile")
	checkFatal(t, err)

	entries, err := statusEntries(repo)
	checkFatal(t, err)

	if len(entries) != 1 {
		t.Errorf("Repo with single unstaged file should have only one status entry")
	}

	if entries[0].Status != git.StatusWtNew {
		t.Errorf("Invalid status, got %d; want %d", entries[0].Status, git.StatusWtNew)
	}

	statuses, err := GetStatus(repo.Workdir())
	checkFatal(t, err)

	if len(statuses) != 1 && statuses[0] != StatusUntrackedFiles {
		t.Errorf("Empty repo should have a single StatusUntrackedFiles status")
	}
}

func TestStatusWithSingleStagedFile(t *testing.T) {
	repo, err := newTestRepo()
	checkFatal(t, err)
	defer cleanupRepo(t, repo)

	err = createFile(repo, "SomeFile")
	checkFatal(t, err)
	err = stageFile(repo, "SomeFile")
	checkFatal(t, err)

	entries, err := statusEntries(repo)
	checkFatal(t, err)

	if len(entries) != 1 {
		t.Errorf("Repo with single staged file should have only one status entry")
	}

	if entries[0].Status != git.StatusIndexNew {
		t.Errorf("Invalid status, got %d; want %d", entries[0].Status, git.StatusIndexNew)
	}

	statuses, err := GetStatus(repo.Workdir())
	checkFatal(t, err)

	if len(statuses) != 1 && statuses[0] != StatusUncommittedChanges {
		t.Errorf("Empty repo should have a single SStatusUncommittedChange status")
	}
}

func TestStatusWithSingleCommit(t *testing.T) {
	repo, err := newTestRepo()
	checkFatal(t, err)
	defer cleanupRepo(t, repo)

	err = createFile(repo, "SomeFile")
	checkFatal(t, err)
	err = stageFile(repo, "SomeFile")
	checkFatal(t, err)
	err = createCommit(repo, "Initial commit")
	checkFatal(t, err)

	entries, err := statusEntries(repo)
	checkFatal(t, err)

	if len(entries) != 0 {
		t.Errorf("Repo with no uncommitted files should have no status entries")
	}

	statuses, err := GetStatus(repo.Workdir())
	checkFatal(t, err)

	if len(statuses) != 1 && statuses[0] != StatusOk {
		t.Errorf("Repo with no uncommitted files should have a single StatusOk status")
	}
}

func TestStatusWithMultipleCommits(t *testing.T) {
	repo, err := newTestRepo()
	checkFatal(t, err)
	defer cleanupRepo(t, repo)

	err = createFile(repo, "SomeFile")
	checkFatal(t, err)
	err = stageFile(repo, "SomeFile")
	checkFatal(t, err)
	err = createCommit(repo, "Initial commit")
	checkFatal(t, err)

	err = createFile(repo, "AnotherFile")
	checkFatal(t, err)
	err = stageFile(repo, "AnotherFile")
	checkFatal(t, err)
	err = createCommit(repo, "Second commit")
	checkFatal(t, err)

	entries, err := statusEntries(repo)
	checkFatal(t, err)

	if len(entries) != 0 {
		t.Errorf("Repo with no uncommitted files should have no status entries")
	}
	statuses, err := GetStatus(repo.Workdir())
	checkFatal(t, err)

	if len(statuses) != 1 && statuses[0] != StatusOk {
		t.Errorf("Repo with no uncommitted files should have a single StatusOk status")
	}
}
