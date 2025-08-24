package test

import (
	"path/filepath"
	"testing"
)

// Repo represents a test repository.
// It embeds testing.T so that any error thrown while creating a test repo will cause a t.Fatal call.
type Repo struct {
	path string
	t    *testing.T
}

// Path returns path to a repository.
func (r *Repo) Path() string {
	return r.path
}

// RepoEmpty creates an empty git repo.
func RepoEmpty(t *testing.T) *Repo {
	t.Helper()

	return RepoEmptyInDir(t, "")
}

// RepoEmptyInDir creates an empty git repo inside a given parent dir.
func RepoEmptyInDir(t *testing.T, parent string) *Repo {
	t.Helper()
	r := &Repo{
		path: TempDir(t, parent),
		t:    t,
	}

	r.init()

	return r
}

// RepoWithUntracked creates a git repo with a single untracked file.
func RepoWithUntracked(t *testing.T) *Repo {
	t.Helper()
	r := RepoEmpty(t)
	r.writeFile("README.md", "I'm a readme file")

	return r
}

// RepoWithStaged creates a git repo with a single staged file.
func RepoWithStaged(t *testing.T) *Repo {
	t.Helper()
	r := RepoEmpty(t)
	r.writeFile("README.md", "I'm a readme file")
	r.stageFile("README.md")

	return r
}

// RepoWithCommit creates a git repo with a single commit.
func RepoWithCommit(t *testing.T) *Repo {
	t.Helper()
	r := RepoEmpty(t)
	r.writeFile("README.md", "I'm a readme file")
	r.stageFile("README.md")
	r.commit("Initial commit")

	return r
}

// RepoWithUncommittedAndUntracked creates a git repo with one staged but uncommitted file and one untracked file.
func RepoWithUncommittedAndUntracked(t *testing.T) *Repo {
	t.Helper()
	r := RepoEmpty(t)
	r.writeFile("README.md", "I'm a readme file")
	r.stageFile("README.md")
	r.commit("Initial commit")
	r.writeFile("README.md", "These changes won't be committed")
	r.writeFile("untracked.txt", "I'm untracked")

	return r
}

// RepoWithBranch creates a git repo with a new branch.
func RepoWithBranch(t *testing.T) *Repo {
	t.Helper()
	r := RepoWithCommit(t)
	r.branch("feature/branch")
	r.checkout("feature/branch")

	return r
}

// RepoWithTag creates a git repo with a new tag.
func RepoWithTag(t *testing.T) *Repo {
	t.Helper()
	r := RepoWithCommit(t)
	r.tag("v0.0.1")
	r.checkout("v0.0.1")

	return r
}

// RepoWithBranchWithUpstream creates a git repo by cloning another repo and checking out a remote branch.
func RepoWithBranchWithUpstream(t *testing.T) *Repo {
	t.Helper()
	origin := RepoWithCommit(t)
	origin.branch("feature/branch")

	r := origin.clone()
	r.checkout("feature/branch")

	return r
}

// RepoWithBranchWithoutUpstream creates a git repo by cloning another repo and checking out a new local branch.
func RepoWithBranchWithoutUpstream(t *testing.T) *Repo {
	t.Helper()
	origin := RepoWithCommit(t)

	r := origin.clone()
	r.branch("feature/branch")
	r.checkout("feature/branch")

	return r
}

// RepoWithBranchAhead creates a git repo with a branch being ahead of a remote branch by 1 commit.
func RepoWithBranchAhead(t *testing.T) *Repo {
	t.Helper()
	origin := RepoWithCommit(t)
	origin.branch("feature/branch")

	r := origin.clone()
	r.checkout("feature/branch")

	r.writeFile("local.new", "local.new")
	r.stageFile("local.new")
	r.commit("local.new")

	return r
}

// RepoWithBranchBehind creates a git repo with a branch being behind a remote branch by 1 commit.
func RepoWithBranchBehind(t *testing.T) *Repo {
	t.Helper()
	origin := RepoWithCommit(t)
	origin.branch("feature/branch")
	origin.checkout("feature/branch")

	r := origin.clone()
	r.checkout("feature/branch")

	origin.writeFile("origin.new", "origin.new")
	origin.stageFile("origin.new")
	origin.commit("origin.new")

	r.fetch()

	return r
}

// RepoWithBranchAheadAndBehind creates a git repo with a branch being 2 commits ahead and 1 behind a remote branch.
func RepoWithBranchAheadAndBehind(t *testing.T) *Repo {
	t.Helper()
	origin := RepoWithCommit(t)
	origin.branch("feature/branch")
	origin.checkout("feature/branch")

	r := origin.clone()
	r.checkout("feature/branch")

	origin.writeFile("origin.new", "origin.new")
	origin.stageFile("origin.new")
	origin.commit("origin.new")

	r.writeFile("local.new", "local.new")
	r.stageFile("local.new")
	r.commit("local.new")

	r.writeFile("local.new2", "local.new2")
	r.stageFile("local.new2")
	r.commit("local.new2")

	r.fetch()

	return r
}

// RepoWithEmptyConfig creates a git repo with empty .git/config file.
func RepoWithEmptyConfig(t *testing.T) *Repo {
	t.Helper()
	r := RepoEmpty(t)
	r.writeFile(filepath.Join(".git", "config"), "")

	return r
}

// RepoWithValidConfig creates a git repo with valid content in .git/config file.
func RepoWithValidConfig(t *testing.T) *Repo {
	t.Helper()
	r := RepoEmpty(t)

	gitconfig := `
	[user]
		name = grdl
	[gitget]
		host = github.com
	`
	r.writeFile(filepath.Join(".git", "config"), gitconfig)

	return r
}
