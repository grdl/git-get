package new

import (
	"testing"
)

func TestRepoClone(t *testing.T) {
	origin := NewRepoWithCommit(t)
	path := NewTempDir(t)

	repo, err := CloneRepo(origin.URL, path, true)
	checkFatal(t, err)

	wt, err := repo.repo.Worktree()
	checkFatal(t, err)

	files, err := wt.Filesystem.ReadDir("")
	checkFatal(t, err)

	if len(files) == 0 {
		t.Errorf("Cloned repo should contain files")
	}
}

func TestRepoEmpty(t *testing.T) {
	repo := NewRepoEmpty(t)

	wt, err := repo.Repo.Worktree()
	checkFatal(t, err)

	status, err := wt.Status()
	checkFatal(t, err)
	if !status.IsClean() {
		t.Errorf("Empty repo should be clean")
	}
}

func TestRepoWithUntrackedFile(t *testing.T) {
	repo := NewRepoWithUntracked(t)

	wt, err := repo.Repo.Worktree()
	checkFatal(t, err)

	status, err := wt.Status()
	checkFatal(t, err)
	if status.IsClean() {
		t.Errorf("Repo with untracked file should not be clean")
	}

	// TODO: remove magic strings
	if !status.IsUntracked("README") {
		t.Errorf("New file should be untracked")
	}
}

func TestRepoWithStagedFile(t *testing.T) {
	repo := NewRepoWithStaged(t)

	wt, err := repo.Repo.Worktree()
	checkFatal(t, err)

	status, err := wt.Status()
	checkFatal(t, err)
	if status.IsClean() {
		t.Errorf("Repo with staged file should not be clean")
	}

	if status.IsUntracked("README") {
		t.Errorf("Staged file should not be untracked")
	}
}

func TestRepoWithSingleCommit(t *testing.T) {
	repo := NewRepoWithCommit(t)

	wt, err := repo.Repo.Worktree()
	checkFatal(t, err)

	status, err := wt.Status()
	checkFatal(t, err)
	if !status.IsClean() {
		t.Errorf("Repo with committed file should be clean")
	}

	if status.IsUntracked("README") {
		t.Errorf("Committed file should not be untracked")
	}
}
func TestStatusWithModifiedFile(t *testing.T) {
	//todo modified but not staged
}

func TestStatusWithUntrackedButIgnoredFile(t *testing.T) {
	//todo
}
