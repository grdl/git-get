package pkg

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

	status, err := loadStatus(repo)
	checkFatal(t, err)

	want := &RepoStatus{
		HasUntrackedFiles:     false,
		HasUncommittedChanges: false,
		Branches:              status.Branches,
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

	status, err := loadStatus(repo)
	checkFatal(t, err)

	want := &RepoStatus{
		HasUntrackedFiles:     true,
		HasUncommittedChanges: false,
		Branches:              status.Branches,
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

	status, err := loadStatus(repo)
	checkFatal(t, err)

	want := &RepoStatus{
		HasUntrackedFiles:     false,
		HasUncommittedChanges: true,
		Branches:              status.Branches,
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

	status, err := loadStatus(repo)
	checkFatal(t, err)

	want := &RepoStatus{
		HasUntrackedFiles:     false,
		HasUncommittedChanges: false,
		Branches:              status.Branches,
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

	status, err := loadStatus(repo)
	checkFatal(t, err)

	want := &RepoStatus{
		HasUntrackedFiles:     false,
		HasUncommittedChanges: false,
		Branches:              status.Branches,
	}

	if !reflect.DeepEqual(status, want) {
		t.Errorf("Wrong repo status, got %+v; want %+v", status, want)
	}
}

func TestStatusCloned(t *testing.T) {
	origin := newTestRepo(t)
	repoRoot := newTempDir(t)

	path, err := CloneRepo(origin.Path(), repoRoot)
	checkFatal(t, err)
	repo, err := OpenRepo(path)
	checkFatal(t, err)

	status, err := loadStatus(repo.repo)
	checkFatal(t, err)

	want := &RepoStatus{
		HasUntrackedFiles:     false,
		HasUncommittedChanges: false,
		Branches:              status.Branches,
	}

	if !reflect.DeepEqual(status, want) {
		t.Errorf("Wrong repo status, got %+v; want %+v", status, want)
	}
}

func TestBranchNewLocal(t *testing.T) {
	repo := newTestRepo(t)

	createFile(t, repo, "file")
	stageFile(t, repo, "file")
	createCommit(t, repo, "Initial commit")
	branch := createBranch(t, repo, "branch")

	status, err := branchStatus(branch)
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

func TestBranchCloned(t *testing.T) {
	origin := newTestRepo(t)
	createFile(t, origin, "file")
	stageFile(t, origin, "file")
	createCommit(t, origin, "Initial commit")

	createBranch(t, origin, "branch")

	repoRoot := newTempDir(t)
	path, err := CloneRepo(origin.Path(), repoRoot)
	checkFatal(t, err)
	repo, err := OpenRepo(path)
	checkFatal(t, err)

	createBranch(t, repo.repo, "local")

	checkoutBranch(t, repo.repo, "branch")
	createFile(t, repo.repo, "anotherFile")
	stageFile(t, repo.repo, "anotherFile")
	createCommit(t, repo.repo, "Second commit")

	err = repo.Reload()
	checkFatal(t, err)

	var tests = []struct {
		got  BranchStatus
		want BranchStatus
	}{
		{repo.Status.Branches["master"], BranchStatus{
			Name:        "master",
			IsRemote:    false,
			HasUpstream: true,
		}},
		{repo.Status.Branches["origin/master"], BranchStatus{
			Name:        "origin/master",
			IsRemote:    true,
			HasUpstream: false,
		}},
		{repo.Status.Branches["branch"], BranchStatus{
			Name:        "branch",
			IsRemote:    false,
			HasUpstream: true,
			Ahead:       1,
			NeedsPush:   true,
		}},
		{repo.Status.Branches["origin/branch"], BranchStatus{
			Name:        "origin/branch",
			IsRemote:    true,
			HasUpstream: false,
		}},
		{repo.Status.Branches["local"], BranchStatus{
			Name:        "local",
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
