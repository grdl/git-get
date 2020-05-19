package main

import (
	"reflect"
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

	status, err := NewRepoStatus(repo.Workdir())
	checkFatal(t, err)

	want := RepoStatus{
		HasUntrackedFiles:     false,
		HasUncommittedChanges: false,
		BranchStatuses:        nil,
	}

	if !reflect.DeepEqual(status, want) {
		t.Errorf("Wrong repo status, got %+v; want %+v", status, want)
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

	status, err := NewRepoStatus(repo.Workdir())
	checkFatal(t, err)

	want := RepoStatus{
		HasUntrackedFiles:     true,
		HasUncommittedChanges: false,
		BranchStatuses:        nil,
	}

	if !reflect.DeepEqual(status, want) {
		t.Errorf("Wrong repo status, got %+v; want %+v", status, want)
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

	status, err := NewRepoStatus(repo.Workdir())
	checkFatal(t, err)

	want := RepoStatus{
		HasUntrackedFiles:     false,
		HasUncommittedChanges: true,
		BranchStatuses:        nil,
	}

	if !reflect.DeepEqual(status, want) {
		t.Errorf("Wrong repo status, got %+v; want %+v", status, want)
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

	status, err := NewRepoStatus(repo.Workdir())
	checkFatal(t, err)

	want := RepoStatus{
		HasUntrackedFiles:     false,
		HasUncommittedChanges: false,
		BranchStatuses:        nil,
	}

	if !reflect.DeepEqual(status, want) {
		t.Errorf("Wrong repo status, got %+v; want %+v", status, want)
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
	status, err := NewRepoStatus(repo.Workdir())
	checkFatal(t, err)

	want := RepoStatus{
		HasUntrackedFiles:     false,
		HasUncommittedChanges: false,
		BranchStatuses:        nil,
	}

	if !reflect.DeepEqual(status, want) {
		t.Errorf("Wrong repo status, got %+v; want %+v", status, want)
	}
}
