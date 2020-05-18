package main

import (
	"testing"

	git "github.com/libgit2/git2go/v30"
)

func TestStatusEntriesWithEmptyRepo(t *testing.T) {
	repo, err := newTestRepo()
	checkFatal(t, err)
	defer cleanupRepo(t, repo)

	entries, err := statusEntries(repo)
	checkFatal(t, err)

	if len(entries) != 0 {
		t.Errorf("Empty repo should have no status entries")
	}
}

func TestStatusEntriesWithSingleUnstagedFile(t *testing.T) {
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
}

func TestStatusEntriesWithSingleStagedFile(t *testing.T) {
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
}

func TestStatusEntriesWithSingleCommit(t *testing.T) {
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
}

func TestStatusEntriesWithMultipleCommits(t *testing.T) {
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
}
