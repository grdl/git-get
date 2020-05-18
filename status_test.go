package main

import (
	"testing"

	git "github.com/libgit2/git2go/v30"
)

func TestStatusWithEmptyRepo(t *testing.T) {
	repo := newTestRepo(t)

	entries, err := statusEntries(repo)
	checkFatal(t, err)

	if len(entries) != 0 {
		t.Errorf("Empty repo should have no status entries")
	}

	status, err := Status(repo.Workdir())
	checkFatal(t, err)

	if status.HasUntrackedFiles != false {
		t.Errorf("Repo should not have untracked files")
	}

	if status.HasUncommittedChanges != false {
		t.Errorf("Repo should not have uncommitted changes")
	}
}

func TestStatusWithUntrackedFile(t *testing.T) {
	repo := newTestRepo(t)
	createFile(t, repo, "SomeFile")

	entries, err := statusEntries(repo)
	checkFatal(t, err)

	if len(entries) != 1 {
		t.Errorf("Repo with untracked file should have only one status entry")
	}

	if entries[0].Status != git.StatusWtNew {
		t.Errorf("Invalid status, got %d; want %d", entries[0].Status, git.StatusWtNew)
	}

	status, err := Status(repo.Workdir())
	checkFatal(t, err)

	if status.HasUntrackedFiles != true {
		t.Errorf("Repo should have untracked files")
	}

	if status.HasUncommittedChanges != false {
		t.Errorf("Repo should not have uncommitted changes")
	}
}

func TestStatusWithUnstagedFile(t *testing.T) {
	//todo
}

func TestStatusWithUntrackedButIgnoredFile(t *testing.T) {
	//todo
}

func TestStatusWithStagedFile(t *testing.T) {
	repo := newTestRepo(t)
	createFile(t, repo, "SomeFile")
	stageFile(t, repo, "SomeFile")

	entries, err := statusEntries(repo)
	checkFatal(t, err)

	if len(entries) != 1 {
		t.Errorf("Repo with staged file should have only one status entry")
	}

	if entries[0].Status != git.StatusIndexNew {
		t.Errorf("Invalid status, got %d; want %d", entries[0].Status, git.StatusIndexNew)
	}

	status, err := Status(repo.Workdir())
	checkFatal(t, err)

	if status.HasUntrackedFiles != false {
		t.Errorf("Repo should not have untracked files")
	}

	if status.HasUncommittedChanges != true {
		t.Errorf("Repo should have uncommitted changes")
	}
}

func TestStatusWithSingleCommit(t *testing.T) {
	repo := newTestRepo(t)
	createFile(t, repo, "SomeFile")
	stageFile(t, repo, "SomeFile")
	createCommit(t, repo, "Initial commit")

	entries, err := statusEntries(repo)
	checkFatal(t, err)

	if len(entries) != 0 {
		t.Errorf("Repo with no uncommitted files should have no status entries")
	}

	status, err := Status(repo.Workdir())
	checkFatal(t, err)

	if status.HasUntrackedFiles != false {
		t.Errorf("Repo should not have untracked files")
	}

	if status.HasUncommittedChanges != false {
		t.Errorf("Repo should not have uncommitted changes")
	}
}

func TestStatusWithMultipleCommits(t *testing.T) {
	repo := newTestRepo(t)
	createFile(t, repo, "SomeFile")
	stageFile(t, repo, "SomeFile")
	createCommit(t, repo, "Initial commit")
	createFile(t, repo, "AnotherFile")
	stageFile(t, repo, "AnotherFile")
	createCommit(t, repo, "Second commit")

	entries, err := statusEntries(repo)
	checkFatal(t, err)

	if len(entries) != 0 {
		t.Errorf("Repo with no uncommitted files should have no status entries")
	}
	status, err := Status(repo.Workdir())
	checkFatal(t, err)

	if status.HasUntrackedFiles != false {
		t.Errorf("Repo should not have untracked files")
	}

	if status.HasUncommittedChanges != false {
		t.Errorf("Repo should not have uncommitted changes")
	}
}
